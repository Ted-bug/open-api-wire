package model

import (
	"api-gin/infra/model"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// HelloWorld 对应数据库 hello_world 表的 GORM 模型
type HelloWorld struct {
	ID         int64                 `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement" json:"id"`
	Name       model.Null[string]    `gorm:"column:name;type:varchar(255)" json:"name"`
	CreateTime model.Null[time.Time] `gorm:"column:create_time;type:datetime" json:"create_time"`
	Age        model.Null[int]       `gorm:"column:age;type:int(11)" json:"age"`
	MyStruct   model.Null[MyStruct]  `gorm:"column:mystruct;type:varchar(255)" json:"mystruct"`
}

// TableName 指定表名
func (h *HelloWorld) TableName() string {
	return "hello_world"
}

type MyStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age,omitempty"`
}

func (m *MyStruct) Scan(value any) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, m)
}

// Value 实现 driver.Valuer 接口
func (m MyStruct) Value() (driver.Value, error) {
	return json.Marshal(m)
}
