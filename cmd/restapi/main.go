package main

//bytes := make([]byte, 32)
//msg := ""
//rand.Read(bytes)
////
////for _, v := range bytes {
////	msg += string(v)
////}
//msg = hex.EncodeToString(bytes)
import (
	"auth-server/config"
	"auth-server/internal/app/presenter/http/handler"
	"auth-server/internal/app/presenter/http/server"
	"auth-server/internal/app/service"
	"auth-server/internal/app/service/user_service"
	ms "auth-server/internal/app/store/mongo_store"
	"auth-server/internal/app/utils/validators"
	"auth-server/pkg/emailsender"
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}
	//Init mongo
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Err in init database. Err message: %s", err.Error())
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Database ping error. %s", err.Error())
	}
	defer client.Disconnect(ctx)
	db := client.Database(viper.GetString("mongo.name"))
	//Init store
	store := ms.NewStore(db)
	//Init repositories
	store.User()
	store.Client()

	uvalidator, err := validators.NewUserValidator()
	if err != nil {
		log.Fatalf("Err in init user validator. Err message: %s", err.Error())
	}
	usvc, err := user_service.New(store, uvalidator)
	if err != nil {
		log.Fatalf("Err in init user service. Err message: %s", err.Error())
	}
	svm, err := service.NewManager(usvc, nil)
	if err != nil {
		log.Fatalf("Err in init user manager. Err message: %s", err.Error())
	}

	email := viper.GetString("email.email")
	password := viper.GetString("email.password")
	host := viper.GetString("email.host")
	port := viper.GetString("port")
	companyName := viper.GetString("email.companyName")

	emailSender := emailsender.New(email, password, host, port, companyName, email)

	userhandler := handler.NewUserHandler(svm, emailSender)

	server, err := server.NewServer(svm, store, userhandler)
	if err != nil {
		log.Fatalf("Error creating server, err: %s", err.Error())
	}
	log.Fatal(server.Run())
}
