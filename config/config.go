package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Name  string       `mapstructure:"name"`
	Host  string       `mapstructure:"host"`
	Port  int          `mapstructure:"port"`
	Mode  string       `mapstructure:"mode"`
	Log   LogrusConfig `mapstructure:"log"`
	MySQL MySQLConfig  `mapstructure:"mysql"`
	Redis RedisConfig  `mapstructure:"redis"`
}

type LogrusConfig struct {
	Level      string `mapstructure:"level"`  // trace, debug, info, warn, error, fatal, panic
	Format     string `mapstructure:"format"` // text, json
	Mode       string `mapstructure:"output"` // command,file
	Path       string `mapstructure:"path"`
	FileName   string `mapstructure:"filename"`
	MaxAge     int    `mapstructure:"max_age"`     // 保留天数，单位天
	MaxSize    int    `mapstructure:"max_size"`    // 保留日志文件大小，单位MB
	MaxBackups uint   `mapstructure:"max_backups"` // 保留份数，单位个
}

type MySQLConfig struct {
	Master          []string `mapstructure:"master"`
	Slave           []string `mapstructure:"slave"`
	Log             string   `mapstructure:"log"` // info, warn, error
	MaxIdleConns    int      `mapstructure:"max_idle_conns"`
	MaxOpenConns    int      `mapstructure:"max_open_conns"`
	ConnMaxLifetime int      `mapstructure:"conn_max_lifetime"`  // 单位 秒
	ConnMaxIdleTime int      `mapstructure:"conn_max_idle_time"` // 单位 秒
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func NewConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置默认值
	viper.SetDefault("logrus.level", "info")
	viper.SetDefault("logrus.format", "text")
	viper.SetDefault("logrus.output", "stdout")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
