package services

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvData struct {
	AccessKey  string
	SecretKey  string
	MongoDbCon string
}

func InitEnvData() EnvData {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	return EnvData{
		AccessKey:  os.Getenv("ACCESSKEYID"),
		SecretKey:  os.Getenv("SECRETACCESSKEY"),
		MongoDbCon: os.Getenv("MONGODBCONNECTION"),
	}
}
