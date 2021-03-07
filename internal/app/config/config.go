package config

import (
	"os"
)

type Config struct {
	MongoURI             string
	MongoDatabase        string
	MongoUsername        string
	MongoPassword        string
	MongoUsersColName    string
	MongoClientColName   string
	AppLink              string
	JWTKey               string
	EmailConfKey         string
	EmailHost            string
	EmailHostPort        string
	CompanyEmail         string
	CompanyEmailPassword string
	CompanyName          string
}

var Cfg = GetConfig()

func GetConfig() *Config {

	return &Config{
		MongoURI:             getEnv("MONGO_URI", ""),
		MongoDatabase:        getEnv("MONGO_DATABASE", "auth"),
		MongoUsername:        getEnv("MONGO_USERNAME", ""),
		MongoPassword:        getEnv("MONGO_PASSWORD", ""),
		MongoUsersColName:    getEnv("MONGO_USERS_COLLECTION", "users"),
		MongoClientColName:   getEnv("MONGO_CLIENTS_COLLECTION", "clients"),
		AppLink:              getEnv("APPLICATION_LINK", ""),
		JWTKey:               getEnv("JWT_KEY", ""),
		EmailConfKey:         getEnv("EMAIL_CONF_KEY", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
		EmailHost:            getEnv("EMAIL_HOST", ""),
		EmailHostPort:        getEnv("EMAIL_HOST_PORT", ""),
		CompanyEmail:         getEnv("COMPANY_EMAIL", "example@examle.org"),
		CompanyEmailPassword: getEnv("COMPANY_EMAIL_PASSWORD", "password"),
		CompanyName:          getEnv("COMPANY_NAME", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

//Init func initialize a viper config from file
//func Init() error {
//	if DEV {
//		viper.SetConfigName("config.dev")
//	} else {
//		viper.SetConfigName("config")
//	}
//	viper.AddConfigPath("config")
//	viper.SetConfigType("yml")
//	err := viper.ReadInConfig()
//	if err != nil {
//		viper.SetConfigName("config.dev")
//	}
//	err = viper.ReadInConfig()
//	return err
//}
