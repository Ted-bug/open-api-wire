package config

import (
	"api-gin/infra/log"
	"api-gin/repo"
	"github.com/spf13/viper"
)

type Config struct {
	Name  string           `mapstructure:"name"`
	Host  string           `mapstructure:"host"`
	Port  int              `mapstructure:"port"`
	Mode  string           `mapstructure:"mode"`
	Log   log.Config       `mapstructure:"log"`
	MySQL repo.MysqlConfig `mapstructure:"mysql"`
	Redis RedisConfig      `mapstructure:"redis"`
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

func GetLogConfig(c *Config) log.Config {
	return c.Log
}

func GetMySQLConfig(c *Config) repo.MysqlConfig {
	return c.MySQL
}

func GetRedisConfig(c *Config) RedisConfig {
	return c.Redis
}
