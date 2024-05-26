package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/Netflix/go-env"
)

var config *Config

type Config struct {
	OpenAIAPIKey                 string  `env:"OPENAI_API_KEY"`
	DatabaseHost                 string  `env:"DATABASE_HOST"`
	DatabaseUser                 string  `env:"DATABASE_USER"`
	DatabasePassword             string  `env:"DATABASE_PASSWORD"`
	DatabaseName                 string  `env:"DATABASE_NAME"`
	LogLevel                     int     `env:"LOG_LEVEL"`
	SMSQueueURL                  string  `env:"SMS_QUEUE_URL"`
	SNSTopicARN                  string  `env:"SNS_TOPIC_ARN"`
	TwilioSID                    string  `env:"TWILIO_SID"`
	TwilioAuthToken              string  `env:"TWILIO_AUTH_TOKEN"`
	TwilioPhoneNumber            string  `env:"TWILIO_PHONE_NUMBER"`
	TwilioStatusCallbackURL      string  `env:"TWILIO_STATUS_CALLBACK_URL"`
	TwilioVerifyServiceSID       string  `env:"TWILIO_VERIFY_SERVICE_SID"`
	ChatModelName                string  `env:"CHAT_MODEL_NAME"`
	ChatModelTemperature         float32 `env:"CHAT_MODEL_TEMPERATURE"`
	ChatModelMaxCompletionTokens int     `env:"CHAT_MODEL_MAX_COMPLETION_TOKENS"`
	ChatModelFrequencyPenalty    float32 `env:"CHAT_MODEL_FREQUENCY_PENALTY"`
}

func New() *Config {
	return &Config{}
}

func Get() *Config {

	if config != nil {
		return config
	}

	var err error
	config := &Config{}

	if _, err = env.UnmarshalFromEnviron(config); err != nil {
		log.Printf("failed to unmarshal environment: %s\n", err)

		return nil
	}

	if err = validateConfig(config); err != nil {
		log.Printf("failed to validate config: %s.  Exiting.", err)

		return nil
	}

	return config
}
func validateConfig(config *Config) error {

	v := reflect.ValueOf(config)

	// Check if it's a pointer and dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if it's a struct type
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("config is not a struct")
	}

	// Now you can safely call NumField
	m := v.NumField()

	for i := 0; i < m; i++ {
		if v.Field(i).String() == "" {
			return fmt.Errorf("%s must not be empty", v.Type().Field(i).Name)
		}
	}

	return nil
}

var DefaultHttpHeaders = map[string]string{
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent",
	"Access-Control-Allow-Methods": "OPTIONS, GET, POST",
}
