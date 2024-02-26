package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
)

func main() {

	log.New("OTP Lambda booting...").Log()

	cfg := config.Get()

	if cfg == nil {
		log.New("Could not load config")
	}

	database := db.Get(cfg)
	handler := &SignupOTPLambdaHandler{}
	handler.Init(database)

	log.New("Lambda ready. Invoking.").Log()
	lambda.Start(handler.HandleRequest)
}
