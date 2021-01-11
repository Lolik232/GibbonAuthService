package main

import (
	"auth-server/config"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func Hande(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.UserAgent())
	fmt.Println(r.RemoteAddr)
}

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	fmt.Print(viper.Get("mongo.username"))
	http.HandleFunc("/", Hande)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
