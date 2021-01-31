package config

import (
	"github.com/spf13/viper"
)

var (
	DEV bool = true
)

//Init func initialize a viper config from file
func Init() error {
	if DEV == true {
		viper.SetConfigName("config.dev")
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath("config")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		viper.SetConfigName("config.dev")
	}
	err = viper.ReadInConfig()
	return err
}
