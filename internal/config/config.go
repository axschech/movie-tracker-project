package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

type Config struct {
	Port string   `mapstructure:"port"`
	DB   DBConfig `mapstructure:"db"`
}

func MakeConfig() (Config, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()
	viper.BindEnv("port", "PORT")
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("db.database", "DB_DATABASE")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("No config file found, using environment variables: %v\n", err)
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
