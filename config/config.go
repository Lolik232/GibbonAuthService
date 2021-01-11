package config

import "github.com/spf13/viper"

//Init func initialize a viper config from file
func Init() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	err := viper.ReadInConfig()
	if err != nil {
		viper.SetConfigName("config.dev")
	}
	err = viper.ReadInConfig()
	return err
}
