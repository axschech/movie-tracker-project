package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"port"`
}

func MakeConfig() (Config, error) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AutomaticEnv()
	viper.BindEnv("port", "PORT")
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
