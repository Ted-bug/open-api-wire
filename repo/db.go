package repo

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type MysqlConfig struct {
	Master          []string `mapstructure:"master"`
	Slave           []string `mapstructure:"slave"`
	Log             string   `mapstructure:"log"` // info, warn, error
	MaxIdleConns    int      `mapstructure:"max_idle_conns"`
	MaxOpenConns    int      `mapstructure:"max_open_conns"`
	ConnMaxLifetime int      `mapstructure:"conn_max_lifetime"`  // 单位 秒
	ConnMaxIdleTime int      `mapstructure:"conn_max_idle_time"` // 单位 秒
}

func NewDB(c MysqlConfig) (*gorm.DB, error) {
	if len(c.Master) == 0 || len(c.Slave) == 0 {
		return nil, fmt.Errorf("no mysql master or slave config")
	}
	logLevel := logger.Silent
	switch c.Log {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	}

	d, err := gorm.Open(mysql.New(mysql.Config{
		DSN: c.Master[0],
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}
	// 主库
	sources := make([]gorm.Dialector, 0)
	for _, s := range c.Master {
		sources = append(sources, mysql.New(mysql.Config{
			DSN: s,
		}))
	}

	// 从库
	replicas := make([]gorm.Dialector, 0)
	for _, s := range c.Slave {
		cfg := mysql.Config{
			DSN: s,
		}
		replicas = append(replicas, mysql.New(cfg))
	}

	err = d.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:           sources,
			Replicas:          replicas,
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true, // 是否在日志中输出 对应的主从信息
		}).
			SetMaxIdleConns(c.MaxIdleConns).
			SetMaxOpenConns(c.MaxOpenConns).
			SetConnMaxIdleTime(time.Duration(c.ConnMaxIdleTime) * time.Second).
			SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second),
	)
	if err != nil {
		return nil, err
	}

	return d, nil
}
