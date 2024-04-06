package log

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	api "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/models"
)

// Log is a struct that represents a map of
// key/value pairs that can be marshalled into JSON.
// A "message" property is always present in the map,
// but additional key/value pairs can be added using
// the Add() method. Additional types, like errors and
// TwilioMessageInfo structs, can be added using the AddError()
// and AddTwilioMessageInfo() methods, respectively.
type Log struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"data,omitempty"`
	Logger  *log.Logger       `json:"-"`
}

// New creates a new Log struct.
// The typical usage is:
// New("Hello, %s", "World").Add("key", "value").Respond(200)
// New("Hello, %s", "World").AddError(err).Log
// s := New("Hello, %s", "World").Write()
func New(format string, a ...interface{}) *Log {

	return &Log{
		Message: fmt.Sprintf(format, a...),
		Fields:  make(map[string]string),
		Logger:  log.Default(),
	}
}

func (r *Log) Add(key string, value string) *Log {
	r.Fields[key] = value
	return r
}

func (r *Log) AddUser(user *models.User) *Log {

	r.Fields["user_id"] = strconv.FormatInt(user.ID, 10)
	r.Fields["phone_number"] = user.PhoneNumber
	r.Fields["firstname"] = user.Firstname
	r.Fields["lastname"] = user.Lastname
	r.Fields["email"] = user.Email
	r.Fields["phone_verified"] = strconv.FormatBool(user.PhoneVerified)
	r.Fields["account_status"] = user.AccountStatus.Name
	r.Fields["account_status_id"] = strconv.FormatInt(user.AccountStatusID, 10)
	r.Fields["nudge_enabled"] = strconv.FormatBool(user.NudgesEnabled())
	r.Fields["provider_code"] = user.ProviderCode

	return r
}

func (r *Log) AddError(err error) *Log {

	if err != nil {
		r.Add("error", err.Error())
	}

	return r
}

func (r *Log) AddAPIProxyRequest(req *events.APIGatewayProxyRequest) *Log {

	for key, value := range req.Headers {
		r.Add(key, value)
	}

	return r
}

func (r *Log) AddHTTPRequest(req *http.Request) *Log {

	return r.Add("method", req.Method).
		Add("host", req.Host).
		Add("url", req.URL.String()).
		Add("user_agent", req.UserAgent()).
		Add("remote_addr", req.RemoteAddr).
		Add("content-type", req.Header.Get("content-type")).
		Add("x-twilio-signature", req.Header.Get("x-twilio-signature"))
}

func (r *Log) AddSmsResponse(resp *api.ApiV2010Message) *Log {
	respJSON, _ := json.Marshal(resp)

	return r.Add("payload", string(respJSON))
}

func (r *Log) AddSQSEvent(event *events.SQSMessage) *Log {

	return r.Add("message_id", event.MessageId).
		Add("md5_of_body", event.Md5OfBody).
		Add("event_source", event.EventSource).
		Add("event_source_arn", event.EventSourceARN).
		Add("aws_region", event.AWSRegion).
		Add("md5_of_attributes", event.Md5OfMessageAttributes).
		Add("receipt_handle", event.ReceiptHandle).
		Add("body", event.Body)
}

func (r *Log) AddMessage(msg *models.Message) *Log {

	// Remarshal our struct and return it
	smsJSON, _ := json.Marshal(msg)

	r.Add("body", msg.Body).
		Add("type", msg.MessageType.Name).
		Add("from", msg.From.PhoneNumber).
		Add("to", msg.To.PhoneNumber).
		Add("reference_id", *msg.ReferenceID).
		Add("payload", string(smsJSON))

	return r
}

func (r *Log) AddTwilioMessageInfo(msg *models.TwilioMessageInfo) *Log {

	// Remarshal our struct and return it
	smsJSON, _ := json.Marshal(msg)

	r.Add("from", msg.From).
		Add("body", msg.Body).
		Add("message_sid", msg.MessageSid).
		Add("account_sid", msg.AccountSid).
		Add("messaging_service_sid", msg.MessagingServiceSid).
		Add("sms_sid", msg.SmsSid).
		Add("error_code", msg.ErrorCode).
		Add("status", msg.SMSStatus).
		Add("from_city", msg.FromCity).
		Add("from_state", msg.FromState).
		Add("from_country", msg.FromCountry).
		Add("to", msg.To).
		Add("to_city", msg.ToCity).
		Add("to_state", msg.ToState).
		Add("to_country", msg.ToCountry).
		Add("payload", string(smsJSON))

	return r
}

func (r *Log) Write() string {
	s, err := json.Marshal(r)

	if err != nil {
		log.Fatalf("Error marshalling response: %s\n", err)
		return r.Message
	}

	return string(s)
}

func (r *Log) Respond(status int) (events.APIGatewayProxyResponse, error) {

	logEntry := r.Write()

	return events.APIGatewayProxyResponse{
		Headers:    config.DefaultHttpHeaders,
		StatusCode: status,
		Body:       logEntry,
	}, nil
}

func (r *Log) Log() {
	logEntry := r.Write()
	log.Println(logEntry)
}

func FormatBool(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}
