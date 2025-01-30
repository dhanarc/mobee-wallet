package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

type Config struct {
	HTTP struct {
		Port                    int           `mapstructure:"port"`
		ReadHeaderTimeout       time.Duration `mapstructure:"read_header_timeout"`
		ReadTimeout             time.Duration `mapstructure:"read_timeout"`
		IdleTimeout             time.Duration `mapstructure:"idle_timeout"`
		WriteTimeout            time.Duration `mapstructure:"write_timeout"`
		GracefulShutdownTimeout time.Duration `mapstructure:"graceful_shutdown_timeout"`
	} `mapstructure:"http"`
	Database struct {
		DSN string `mapstructure:"dsn"`
	} `mapstructure:"database"`
	Auth struct {
		Key string `mapstructure:"key"`
	}
}

func InitDefaultConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil
	}
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil
	}
	return &cfg
}
