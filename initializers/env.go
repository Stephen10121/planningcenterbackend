package initializers

import (
	"log"
	"os"
)

var Password string
var Credentials string
var ResourcesJsonFilePath string

func SetupEnv() {
	calPassword := os.Getenv("CAL_PASSWORD")

	if len(calPassword) != 0 {
		Password = calPassword
	} else {
		log.Fatalln("CAL_PASSWORD env variable is not set.")
	}

	cred := os.Getenv("CREDENTIALS")

	if len(cred) != 0 {
		Credentials = cred
	} else {
		log.Fatalln("CREDENTIALS env variable is not set.")
	}

	resFile := os.Getenv("RESOURCESFILE")
	if len(resFile) != 0 {
		ResourcesJsonFilePath = resFile
	} else {
		log.Fatalln("RESOURCESFILE env variable is not set.")
	}
}
