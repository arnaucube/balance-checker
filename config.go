package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Currencies map[string]float64
}

var C Config

func MustRead(configPath string) error {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	err := viper.Unmarshal(&C)
	if err != nil {
		return err
	}
	return nil
}
