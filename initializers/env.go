package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

var Password string
var Credentials string

func SetupEnv() {
	envFile, _ := godotenv.Read(".env")

	calPassword := envFile["CAL_PASSWORD"]

	if len(calPassword) != 0 {
		Password = calPassword
	} else {
		log.Fatalln("CAL_PASSWORD env variable is not set.")
	}

	cred := envFile["CREDENTIALS"]

	if len(cred) != 0 {
		Credentials = cred
	} else {
		log.Fatalln("CREDENTIALS env variable is not set.")
	}
}
