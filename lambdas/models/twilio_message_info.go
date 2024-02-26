package models

type TwilioMessageStatus string

const (
	TwilioMessageStatusFailed      TwilioMessageStatus = "failed"
	TwilioMessageStatusDelivered   TwilioMessageStatus = "delivered"
	TwilioMessageStatusUndelivered TwilioMessageStatus = "undelivered"
	TwilioMessageStatusAccepted    TwilioMessageStatus = "accepted"
	TwilioMessageStatusQueued      TwilioMessageStatus = "queued"
	TwilioMessageStatusSending     TwilioMessageStatus = "sending"
	TwilioMessageStatusSent        TwilioMessageStatus = "sent"
	TwilioMessageStatusReceiving   TwilioMessageStatus = "receiving"
	TwilioMessageStatusReceived    TwilioMessageStatus = "received"
	TwilioMessageStatusRead        TwilioMessageStatus = "read" // WhatsApp only
)

// TwilioMessageInfo represents the data sent by Twilio for an SMS callback.
type TwilioMessageInfo struct {
	MessageSid          string `form:"MessageSid"`
	SmsMessageSid       string `form:"SmsMessageSid"`
	SmsSid              string `form:"SmsSid"`
	SMSStatus           string `form:"SmsStatus"`
	AccountSid          string `form:"AccountSid"`
	MessagingServiceSid string `form:"MessagingServiceSid"`
	ErrorCode           string `form:"ErrorCode"`
	ErrorMessage        string `form:"ErrorMessage"`
	MessageStatus       string `form:"MessageStatus"`
	RawDlrDoneDate      string `form:"rawDlrDoneDate"`
	NumMedia            string `form:"NumMedia"`
	NumSegments         string `form:"NumSegments"`
	ApiVersion          string `form:"ApiVersion"`
	Body                string `form:"Body"`
	From                string `form:"From"`
	FromZip             string `form:"FromZip"`
	FromCity            string `form:"FromCity"`
	FromState           string `form:"FromState"`
	FromCountry         string `form:"FromCountry"`
	To                  string `form:"To"`
	ToZip               string `form:"ToZip"`
	ToCity              string `form:"ToCity"`
	ToState             string `form:"ToState"`
	ToCountry           string `form:"ToCountry"`
}

func (t *TwilioMessageInfo) GetTwilioMessageStatus() TwilioMessageStatus {
	return TwilioMessageStatus(t.MessageStatus)
}

func ConvertTwilioStatusToMessageStatus(status TwilioMessageStatus) MessageStatus {
	switch status {
	case TwilioMessageStatusAccepted:
		return NewMessageStatusAccepted()

	case TwilioMessageStatusDelivered:
		return NewMessageStatusDelivered()

	case TwilioMessageStatusFailed:
		return NewMessageStatusFailed()

	case TwilioMessageStatusQueued:
		return NewMessageStatusQueued()

	case TwilioMessageStatusReceiving:
		return NewMessageStatusReceiving()

	case TwilioMessageStatusReceived:
		return NewMessageStatusReceived()

	case TwilioMessageStatusRead:
		return NewMessageStatusRead()

	case TwilioMessageStatusSending:
		return NewMessageStatusSending()

	case TwilioMessageStatusSent:
		return NewMessageStatusSent()

	case TwilioMessageStatusUndelivered:
		return NewMessageStatusFailed()
	}

	return NewMessageStatusUnknown()
}
