package main

import (
	"auth-server/config"
	"fmt"
	"github.com/spf13/viper"
	"regexp"
)

func main() {
	//bytes := make([]byte, 32)
	//msg := ""
	//rand.Read(bytes)
	////
	////for _, v := range bytes {
	////	msg += string(v)
	////}
	//msg = hex.EncodeToString(bytes)
	err := config.Init()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.Get("mongo.username"))
}
