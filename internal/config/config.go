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

type Source struct {
	BaseURL string `mapstructure:"base_url"`
	APIKey  string `mapstructure:"api_key"`
	PIN     string `mapstructure:"pin"`
}

type Config struct {
	Port     string   `mapstructure:"port"`
	DB       DBConfig `mapstructure:"db"`
	TVSource Source   `mapstructure:"tv_source"`
}

func MakeConfig() (Config, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()
	// not sure how to get things working without this
	viper.BindEnv("port", "PORT")
	viper.BindEnv("db.host", "DB_HOST")
	viper.BindEnv("db.user", "DB_USER")
	viper.BindEnv("db.password", "DB_PASSWORD")
	viper.BindEnv("db.database", "DB_DATABASE")
	viper.BindEnv("tv_source.base_url", "TV_SOURCE_BASE_URL")
	viper.BindEnv("tv_source.api_key", "TV_SOURCE_API_KEY")
	viper.BindEnv("tv_source.pin", "TV_SOURCE_PIN")

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
