package main

import (
	cfg "auth-server/internal/app/config"
	"auth-server/internal/app/presenter/http/handler"
	"auth-server/internal/app/presenter/http/server"
	"auth-server/internal/app/service/services"
	ms "auth-server/internal/app/store/mongo_store"
	"auth-server/internal/app/utils/validators"
	"auth-server/pkg/emailsender"
	"context"
	"github.com/subosito/gotenv"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//bytes := make([]byte, 32)
//msg := ""
//rand.Read(bytes)
////
////for _, v := range bytes {
////	msg += string(v)
////}
//msg = hex.EncodeToString(bytes)

func init() {
	if err := gotenv.Load(".env"); err != nil {
		log.Println(".env file not found.")
	}
}

func main() {
	config := cfg.GetConfig()
	//Init mongo
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
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
	db := client.Database(config.MongoDatabase)

	//Init store
	store := ms.NewStore(db)

	//Init repositories
	store.User()
	store.Client()

	uvalidator, err := validators.NewUserValidator()
	if err != nil {
		log.Fatalf("Err in init user validator. Err message: %s", err.Error())
	}

	svm, err := services.NewManager(store, uvalidator)
	if err != nil {
		log.Fatalf("Err in init user manager. Err message: %s", err.Error())
	}

	emailSender := emailsender.New(
		config.CompanyEmail,
		config.CompanyEmailPassword,
		config.EmailHost,
		config.EmailHostPort,
		config.CompanyName,
		config.CompanyEmail,
	)

	userHandler := handler.NewUserHandler(svm, emailSender)

	server, err := server.NewServer(svm, store, userHandler)
	if err != nil {
		log.Fatalf("Error creating server, err: %s", err.Error())
	}
	log.Fatal(server.Run())
}
