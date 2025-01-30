package initializers

import (
	"log"
	"os"
)

var Password string
var Credentials string
var ResendKey string
var PlanningCenterClientId string
var PlanningCenterClientSecret string

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

	resend := os.Getenv("RESEND_KEY")

	if len(cred) != 0 {
		ResendKey = resend
	} else {
		log.Fatalln("RESEND_KEY env variable is not set.")
	}

	clientId := os.Getenv("PLANNINGCENTER_CLIENT_ID")

	if len(cred) != 0 {
		PlanningCenterClientId = clientId
	} else {
		log.Fatalln("PLANNINGCENTER_CLIENT_ID env variable is not set.")
	}

	clientSecret := os.Getenv("PLANNINGCENTER_CLIENT_SECRET")

	if len(cred) != 0 {
		PlanningCenterClientSecret = clientSecret
	} else {
		log.Fatalln("PLANNINGCENTER_CLIENT_SECRET env variable is not set.")
	}
}
