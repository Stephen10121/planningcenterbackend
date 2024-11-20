package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

var Password string

func SetupEnv() {
	envFile, _ := godotenv.Read(".env")

	value := envFile["CAL_PASSWORD"]

	if len(value) != 0 {
		Password = value
	} else {
		log.Fatalln("CAL_PASSWORD env variable is not set.")
	}
}
