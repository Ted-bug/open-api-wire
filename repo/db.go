package repo

import (
	"api-gin/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func NewDB(conf *config.Config) (*gorm.DB, error) {
	if conf == nil {
		return nil, fmt.Errorf("no config")
	}
	mysqlConfig := conf.MySQL
	if len(mysqlConfig.Master) == 0 || len(mysqlConfig.Slave) == 0 {
		return nil, fmt.Errorf("no mysql master or slave config")
	}
	logLevel := logger.Silent
	switch mysqlConfig.Log {
	case "info":
		logLevel = logger.Info
	case "warn":
		logLevel = logger.Warn
	case "error":
		logLevel = logger.Error
	}

	d, err := gorm.Open(mysql.New(mysql.Config{
		DSN: mysqlConfig.Master[0],
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}
	// 主库
	sources := make([]gorm.Dialector, 0)
	for _, s := range mysqlConfig.Master {
		sources = append(sources, mysql.New(mysql.Config{
			DSN: s,
		}))
	}

	// 从库
	replicas := make([]gorm.Dialector, 0)
	for _, s := range mysqlConfig.Slave {
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
			SetMaxIdleConns(mysqlConfig.MaxIdleConns).
			SetMaxOpenConns(mysqlConfig.MaxOpenConns).
			SetConnMaxIdleTime(time.Duration(mysqlConfig.ConnMaxIdleTime) * time.Second).
			SetConnMaxLifetime(time.Duration(mysqlConfig.ConnMaxLifetime) * time.Second),
	)
	if err != nil {
		return nil, err
	}

	return d, nil
}
