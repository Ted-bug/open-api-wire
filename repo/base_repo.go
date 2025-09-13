package repo

import (
	"context"
	"fmt"
	"gorm.io/plugin/dbresolver"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type IBaseRepo interface {
	// Write 指定写库
	Write(ctx context.Context) *gorm.DB
	// Read 指定读库
	Read(ctx context.Context) *gorm.DB
	// StartTrans 外部调用，启动事务，存入ctx
	StartTrans(ctx context.Context, tx *gorm.DB) context.Context
	// GetDB 每次查询前调用，检查ctx，返回DB，防止启用了事务
	GetDB(ctx context.Context) *gorm.DB
	// SetSuffix 设置分表后缀，参数依据对应的分表实现
	SetSuffix(ctx context.Context, params ...any) (context.Context, error)
	// GetTableName 每次查询前调用，检查ctx，返回表名，防止设置了分表
	GetTableName(ctx context.Context) string
}

type BaseRepo struct {
	Db *gorm.DB

	TableName string
	ShardType ShardType
	ShardFunc *ShardTool
}

// NewBaseRepo 创建一个基础DB，不参与wire
func NewBaseRepo(db *gorm.DB, options ...Option) *BaseRepo {
	baseDB := &BaseRepo{
		Db:        db,
		ShardType: ShardNone,
	}
	for _, option := range options {
		option(baseDB)
	}
	if baseDB.ShardFunc == nil {
		baseDB.ShardFunc = &defaultShardTool
	}
	return baseDB
}

type Option func(*BaseRepo)

func WithTableName(table string) Option {
	return func(repo *BaseRepo) {
		repo.TableName = table
	}
}

// KeyTransDb 事务
type KeyTransDb struct{}

// KeyTableSuffix 分表后缀
type KeyTableSuffix struct{}

// ShardType 分表类型
type ShardType int

const (
	ShardNone ShardType = iota
	ShardTypeDay
	ShardTypeWeek
	ShardTypeMonth
	ShardTypeMod
)

func WithShard(shard bool, t ShardType) Option {
	return func(repo *BaseRepo) {
		repo.ShardType = t
	}
}

type ShardTool struct {
	ShardTypeDay   func(ctx context.Context, t time.Time) context.Context
	ShardTypeWeek  func(ctx context.Context, t time.Time) context.Context
	ShardTypeMonth func(ctx context.Context, t time.Time) context.Context
	ShardTypeMod   func(ctx context.Context, now int, modBase int) context.Context
}

var defaultShardTool = ShardTool{
	ShardTypeDay: func(ctx context.Context, t time.Time) context.Context {
		suffix := t.Format("20060102")
		return context.WithValue(ctx, KeyTableSuffix{}, suffix)
	},
	ShardTypeWeek: func(ctx context.Context, t time.Time) context.Context {
		w := t.Weekday()
		if w == 0 {
			w = 6
		} else {
			w = w - 1
		}
		suffix := t.AddDate(0, 0, -1*int(w)).Format("20060102") + "w"
		return context.WithValue(ctx, KeyTableSuffix{}, suffix)
	},
	ShardTypeMonth: func(ctx context.Context, t time.Time) context.Context {
		suffix := t.Format("200601")
		return context.WithValue(ctx, KeyTableSuffix{}, suffix)
	},
	ShardTypeMod: func(ctx context.Context, now int, modBase int) context.Context {
		suffix := strconv.Itoa(now % modBase)
		return context.WithValue(ctx, KeyTableSuffix{}, suffix)
	},
}

func WithShardFunc(tool *ShardTool) Option {
	return func(repo *BaseRepo) {
		repo.ShardFunc = tool
	}
}

func (b *BaseRepo) Write(ctx context.Context) *gorm.DB {
	return b.Db.WithContext(ctx).Clauses(dbresolver.Write)
}

func (b *BaseRepo) Read(ctx context.Context) *gorm.DB {
	return b.Db.WithContext(ctx).Clauses(dbresolver.Read)
}

func (b *BaseRepo) StartTrans(ctx context.Context, tx *gorm.DB) context.Context {
	if tx == nil {
		tx = b.Write(ctx).Begin()
	}
	return context.WithValue(ctx, KeyTransDb{}, tx)
}

func (b *BaseRepo) GetDB(ctx context.Context) *gorm.DB {
	if txAny := ctx.Value(KeyTransDb{}); txAny != nil {
		if tx, ok := txAny.(*gorm.DB); ok {
			return tx
		}
	}
	return b.Db.WithContext(ctx)
}

func (b *BaseRepo) SetSuffix(ctx context.Context, params ...any) (context.Context, error) {
	switch b.ShardType {
	case ShardNone:
		return ctx, nil
	case ShardTypeDay, ShardTypeWeek, ShardTypeMonth:
		if len(params) == 0 {
			return nil, fmt.Errorf("请提供时间依据")
		}
		if t, ok := params[0].(time.Time); !ok {
			return nil, fmt.Errorf("请提供正确时间")
		} else {
			return b.ShardFunc.ShardTypeDay(ctx, t), nil
		}
	case ShardTypeMod:
		if len(params) != 2 {
			return nil, fmt.Errorf("请提供分表的取模依据")
		}
		now, nowOk := params[0].(int)
		modBase, modBaseOk := params[1].(int)
		if !nowOk || !modBaseOk {
			return nil, fmt.Errorf("请提供正确取模参数")
		}
		return b.ShardFunc.ShardTypeMod(ctx, now, modBase), nil
	default:
		return nil, fmt.Errorf("分表类型错误")
	}
}

func (b *BaseRepo) GetTableName(ctx context.Context) string {
	if b.ShardType == ShardNone {
		return b.TableName
	}
	return fmt.Sprintf("%s_%s", b.TableName, ctx.Value(KeyTableSuffix{}))
}
