package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

var Password string
var AppId string
var AppSecret string

func SetupEnv() {
	envFile, _ := godotenv.Read(".env")

	calPassword := envFile["CAL_PASSWORD"]

	if len(calPassword) != 0 {
		Password = calPassword
	} else {
		log.Fatalln("CAL_PASSWORD env variable is not set.")
	}

	appId := envFile["APP_ID"]

	if len(appId) != 0 {
		AppId = appId
	} else {
		log.Fatalln("APP_ID env variable is not set.")
	}

	appSecret := envFile["APP_SECRET"]

	if len(appSecret) != 0 {
		AppSecret = appSecret
	} else {
		log.Fatalln("APP_SECRET env variable is not set.")
	}
}
