package sqs

import (
	"github.com/kmesiab/equilibria/lambdas/models"
)

type SenderInterface interface {
	Send(queueURL string, message *models.Message) error
}

func NewSQSSender() SenderInterface {
	return &AWSSender{}
}
