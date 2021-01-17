package main

import (
	"auth-server/config"
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.Get("mongo.username"))
}
