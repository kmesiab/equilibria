package twilio

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/twilio/twilio-go"
	"github.com/twilio/twilio-go/client"
	api "github.com/twilio/twilio-go/rest/api/v2010"
	verify "github.com/twilio/twilio-go/rest/verify/v2"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/log"
)

const (
	WebhookContentType = "application/x-www-form-urlencoded"
)

func SendSMS(fromPhoneNumber, toPhoneNumber, message string) (*api.ApiV2010Message, error) {

	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.Get().TwilioSID,
		Password: config.Get().TwilioAuthToken,
	})

	params := &api.CreateMessageParams{
		StatusCallback: &config.Get().TwilioStatusCallbackURL,
	}
	params.SetBody(message)
	params.SetFrom(fromPhoneNumber)
	params.SetTo(toPhoneNumber)

	return twilioClient.Api.CreateMessage(params)

}

func SendOTP(phoneNumber string) (*verify.VerifyV2Verification, error) {

	const (
		channel = "sms"
		locale  = "en"
	)

	twilioRestClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.Get().TwilioSID,
		Password: config.Get().TwilioAuthToken,
	})

	params := &verify.CreateVerificationParams{
		To: &phoneNumber,
	}

	params.SetChannel(channel)
	params.SetLocale(locale)

	response, err := twilioRestClient.VerifyV2.CreateVerification(config.Get().TwilioVerifyServiceSID, params)

	if err != nil {
		return nil, err
	}

	return response, nil

}

func VerifyOTP(phoneNumber string, code string) (*verify.VerifyV2VerificationCheck, error) {

	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: config.Get().TwilioSID,
		Password: config.Get().TwilioAuthToken,
	})

	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := twilioClient.VerifyV2.CreateVerificationCheck(config.Get().TwilioVerifyServiceSID, params)

	if err != nil {
		return nil, err
	}

	return resp, nil

}

func IsValidPhoneNumber(phoneNumber string) bool {
	e164Regex := `^\+[1-9]\d{1,14}$`
	re := regexp.MustCompile(e164Regex)
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")

	return re.Find([]byte(phoneNumber)) != nil
}

// IsValidWebhookRequest validates the Twilio webhook request.  If simpleAuth is true,
// then the request is just checked for the presence of the `X-Twilio-Signature` header.
// the
func IsValidWebhookRequest(request events.APIGatewayProxyRequest, twilioAuthToken string, simpleAuth bool) bool {
	fullURL := BuildRequestURLFromProxyRequest(request)

	// The header is sometimes capitalized, so we call
	// `GetTwilioSignatureFromHeaders()` to range the headers in lower case
	twilioSignature := GetTwilioSignatureFromHeaders(request.Headers)

	// If `GetTwilioSignatureFromHeaders()` returns an empty string, then
	// the signature was not found in the request headers.
	if twilioSignature == "" {
		log.New("No Twilio signature found in request headers").
			AddAPIProxyRequest(&request).Log()

		return false
	}

	formValues, err := url.ParseQuery(request.Body)

	if err != nil {
		log.New("Error parsing request body").
			AddAPIProxyRequest(&request).
			AddError(err).
			Log()

		return false
	}

	paramsMap := ConvertURLValuesToMap(formValues)

	// If `simpleAuth` is true, then we just check for the presence of the
	// `X-Twilio-Signature` header and a match of the correct SmsSid.
	if simpleAuth {

		sidMatch := config.Get().TwilioSID == paramsMap["SmsSid"]
		sigMatch := twilioSignature != ""

		return sigMatch && sidMatch
	}

	requestValidator := client.NewRequestValidator(twilioAuthToken)

	l := log.New("Validating parameters.").
		Add("twilio_sid", config.Get().TwilioSID).
		Add("full_url", fullURL).
		Add("twilio_signature", twilioSignature).
		Add("twilio_auth_token", twilioAuthToken).
		AddAPIProxyRequest(&request).
		AddError(err)

	for key, value := range paramsMap {
		l.Add(key, value)
	}

	l.Log()

	// return requestValidator.Validate(fullURL, paramsMap, twilioSignature)
	return requestValidator.ValidateBody(fullURL, []byte(request.Body), twilioSignature)
}

func BuildRequestURLFromProxyRequest(request events.APIGatewayProxyRequest) string {

	scheme := "https://" // Default scheme
	if forwardedProto, ok := request.Headers["X-Forwarded-Proto"]; ok {
		scheme = forwardedProto + "://"
	}

	host := request.Headers["Host"]
	path := request.Path

	return scheme + host + "/dev" + path // + queryString
}

func ConvertURLValuesToMap(values url.Values) map[string]string {
	result := make(map[string]string)
	for key, valueArray := range values {
		if len(valueArray) > 0 {
			result[key] = valueArray[0]
		}
	}
	return result
}

func GetTwilioSignatureFromHeaders(headers map[string]string) string {

	for key, value := range headers {
		if strings.ToLower(key) == "x-twilio-signature" {

			return value
		}
	}

	return ""
}
