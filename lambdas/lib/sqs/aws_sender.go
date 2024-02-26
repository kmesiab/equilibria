package sqs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/kmesiab/equilibria/lambdas/models"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

type AWSSender struct{}

func (s *AWSSender) Send(queueURL string, message *models.Message) error {

	jsonBody, err := json.Marshal(message)
	bodyString := string(jsonBody)

	if err != nil {
		return err
	}

	// Load AWS SDK configuration from the shared config file (~/.aws/config)
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	// Transmit to the SQS queue
	var sqsClient = sqs.NewFromConfig(cfg)

	result, err := sqsClient.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    &queueURL,
		MessageBody: &bodyString,
	})

	if err != nil {
		return err
	}

	message.ReferenceID = result.MessageId

	return nil

}
