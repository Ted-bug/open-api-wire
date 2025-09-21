package model

import (
	"api-gin/infra/model"
	"api-gin/repo"
	"context"
	"testing"
)

func TestChangeNullValue(t *testing.T) {
	db, err := repo.NewDB(repo.MysqlConfig{
		Master:          []string{"root:root@tcp(localhost:3306)/ai_hzc_agent?charset=utf8mb4&parseTime=true&loc=Local"},
		Slave:           []string{"root:root@tcp(localhost:3306)/ai_hzc_agent?charset=utf8mb4&parseTime=true&loc=Local"},
		Log:             "info",
		MaxIdleConns:    5,
		MaxOpenConns:    50,
		ConnMaxLifetime: 3600,
		ConnMaxIdleTime: 1800,
	})
	if err != nil {
		t.Errorf("Error creating database: %v", err)
	}

	ctx := context.Background()

	// 新增
	//h := HelloWorld{
	//	Name: model.NewNullValid("hello"),
	//	Age:  model.NewNullValid(18),
	//}
	//if err := db.WithContext(ctx).Create(&h).Error; err != nil {
	//	t.Errorf("Error creating record: %v", err)
	//}

	// 查询
	//where := HelloWorld{
	//	Name: model.NewNullValid("hello"),
	//}
	//var h HelloWorld
	//if err := db.WithContext(ctx).Where(&where).First(&h).Error; err != nil {
	//	t.Errorf("Error search record: %v", err)
	//}
	//t.Log(h)
	//jsonData, _ := json.Marshal(h)
	//t.Log(string(jsonData))

	// 设置字段为null
	// sql.Null[T]：V != 零值 且 Valid = false 时， gorm会set Field = Null
	updateData := HelloWorld{
		ID: 3,
		MyStruct: model.NewNullValid(MyStruct{
			Name: "12312",
		}),
	}
	tx := db.WithContext(ctx)
	//tx.DryRun = true
	if err := tx.Updates(&updateData).Error; err != nil {
		t.Errorf("Error update record: %v", err)
	}
	t.Log(updateData)
}
