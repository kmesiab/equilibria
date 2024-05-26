package sqs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/kmesiab/equilibria/lambdas/lib/log"
	"github.com/kmesiab/equilibria/lambdas/models"
)

type SNSSender struct {
	client *sns.Client
}

func NewSNSSender() (*SNSSender, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := sns.NewFromConfig(cfg)
	return &SNSSender{client: client}, nil
}

func (s *SNSSender) Send(topicARN string, message *models.Message) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(messageJSON)),
		TopicArn: aws.String(topicARN),
	}

	_, err = s.client.Publish(context.TODO(), input)

	if err != nil {
		log.New("Publish to SNS error: %s", err).Log()
		return err
	}

	return nil
}
