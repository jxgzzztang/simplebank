package util

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"time"
)

type JWT struct {
	SecretKey string `mapstructure:"SECRET_KEY"`
	ExpiresDuration time.Duration `mapstructure:"EXPIRE_TIME"`
	RefreshDuration time.Duration `mapstructure:"REFRESH_DURATION"`
	Issuer string `mapstructure:"ISSUER"`
}

type ViperConfig struct {
	DBSource string `mapstructure:"dbSource"`
	Port     string `mapstructure:"port"`
	Jwt      JWT `mapstructure:"jwt"`
}

var Config ViperConfig

func LoadConfig(path string) (err error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// Config file not found; ignore error if desired
			fmt.Println("No config file found")
			return err
		}
		fmt.Println("Error loading config file:", err.Error())
		return err
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		return err
	}
	return nil
}