package test

import (
	"database/sql"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aws/aws-lambda-go/events"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/kmesiab/equilibria/lambdas/lib/config"
	"github.com/kmesiab/equilibria/lambdas/lib/db"
	"github.com/kmesiab/equilibria/lambdas/models"
)

// Shared SQL
const (
	ConversationSelectQuery = "SELECT \\* FROM `conversations`"
	MessageSelectQuery      = "SELECT \\* FROM `messages`"
)

// Conversation/Message constants
const (
	DefaultTestReferenceID    = "SMa74e33ba8361485"
	DefaultTestConversationID = 1
	DefaultTestMessage        = "Hello World"
)

// User type constants
const (
	DefaultFromUserID        = 1
	DefaultTestUserFirstname = "John"
	DefaultTestUserLastname  = "Doe"
	DefaultTestPhoneNumber   = "+12533243071"
	DefaultTestEmail         = "john.doe@example.com"
	DefaultTestPassword      = "password123"
)

// JsonError is a common error structure returned by some APIs
// This convenience struct lets us easily unmarshal one, or return one.
type JsonError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// GenerateMockTwilioFormPostRequest are Used in Receive and Send SMS Lambdas
func GenerateMockTwilioFormPostRequest() *events.APIGatewayProxyRequest {

	formData := url.Values{}
	formData.Set("AccountSid", "AC0e8e16c274b3ae7740b1a854b8c9846a")
	formData.Set("ApiVersion", "2010-04-01")
	formData.Set("SmsMessageSid", "SMa74e33ba8361485b4bfbb6ec285ceac5")
	formData.Set("MessageSid", "SMa74e33ba8361485b4bfbb6ec285ceac5")
	formData.Set("SmsSid", "SMa74e33ba8361485b4bfbb6ec285ceac5")
	formData.Set("MessagingServiceSid", "MGa3799c565299f143097ff388571be2b2")

	formData.Set("From", "+12533243071")
	formData.Set("FromCountry", "US")
	formData.Set("FromState", "WA")
	formData.Set("FromCity", "SEATTLE")
	formData.Set("FromZip", "98106")

	formData.Set("To", "+18333595081")
	formData.Set("ToCountry", "US")
	formData.Set("ToState", "")
	formData.Set("ToCity", "")
	formData.Set("ToZip", "")

	formData.Set("Body", "Will I win the lotto")
	formData.Set("SmsStatus", "received")
	formData.Set("MessageStatus", "received")

	formData.Set("NumMedia", "0")
	formData.Set("NumSegments", "1")

	return &events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Headers: map[string]string{
			"Content-Type":       "application/x-www-form-urlencoded",
			"x-twilio-signature": "s2+/wx/g0ehm1dV+s0lCrG/Jn+k=",
			"Host":               "foo.com",
		},
		Body: formData.Encode(),
	}

}

func GenerateMockAccountStatusPending() *sqlmock.Rows {
	testAccountStatusesColumns := sqlmock.NewRows(
		[]string{"id", "name"},
	)

	return testAccountStatusesColumns.AddRow(1, "Pending Activation")
}

func GenerateMockAccountStatusActive() *sqlmock.Rows {
	testAccountStatusesColumns := sqlmock.NewRows(
		[]string{"id", "name"},
	)

	return testAccountStatusesColumns.AddRow(2, "Active")
}

// GenerateMockUserRepositoryUser generates a mocked row for the user table
// and should be used with database mocks where the expected return will be
// a single user row. You may call this multiple times to create several users,
// though they will have the same ID.
func GenerateMockUserRepositoryUser() *sqlmock.Rows {
	columns := sqlmock.NewRows(
		[]string{
			"id",
			"phone_number",
			"phone_verified",
			"firstname",
			"lastname",
			"email",
			"account_status_id",
			"provider_code",
			"nudge_enabled",
			"created_at",
			"deleted_at",
			"updated_at",
		})

	now := time.Now() // Assuming 'now' is the current time for the test

	// Mock data row for messages
	rows := columns.AddRow(
		1,                   // id
		"2533243071",        // phone_number
		true,                // phone_verified
		"jane",              // firstname
		"doe",               // lastname
		"janedoe@email.com", // email
		1,                   // account_status_id
		"CODE",              // provider_code
		true,                // nudge_enabled
		now,                 // created_at
		now,                 // updated_at
		nil,                 // deleted_at
	)

	return rows
}

func ExpectMockSelectMessageStatusAndTypes(mock *sqlmock.Sqlmock) {
	(*mock).ExpectQuery("SELECT \\* FROM `message_statuses`").
		WithArgs(models.NewMessageStatusPending().ID).WillReturnRows(GenerateMockMessageStatus())

	(*mock).ExpectQuery("SELECT \\* FROM `message_types`").
		WithArgs(models.NewMessageTypeSMS().ID).WillReturnRows(GenerateMockMessageType())
}

func ExpectMockSelectUser(mock *sqlmock.Sqlmock, arg interface{}) {
	(*mock).ExpectQuery("SELECT \\* FROM `users`").
		WithArgs(arg).WillReturnRows(GenerateMockUserRepositoryUser())
	(*mock).ExpectQuery("SELECT \\* FROM `account_statuses`").
		WithArgs(1).WillReturnRows(GenerateMockAccountStatusPending())
}

func ExpectMockInsertMessage(mock *sqlmock.Sqlmock) {
	(*mock).ExpectBegin()
	(*mock).ExpectExec("INSERT INTO `messages`").
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnResult(GenerateMockLastAffectedRow())
	(*mock).ExpectCommit()
}

func ExpectMockInsertConversation(mock *sqlmock.Sqlmock) {
	(*mock).ExpectBegin()
	(*mock).ExpectExec("INSERT INTO `conversations`").WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnResult(GenerateMockLastAffectedRow())
	(*mock).ExpectCommit()
}

// SQL Mock Types, like column names, mock records from
// select and update statements.

func GenerateMockMessageType() *sqlmock.Rows {
	return sqlmock.NewRows(
		[]string{"id", "name", "bill_rate_in_credits"},
	).AddRow(2, "SMS", 0.5)
}

func GenerateMockMessageRepositoryMessages() *sqlmock.Rows {
	// Defining the columns for the message table based on your schema
	testMessagesColumns := sqlmock.NewRows(
		[]string{
			"id",
			"reference_id",
			"conversation_id",
			"from_user_id",
			"to_user_id",
			"body",
			"message_type_id",
			"sent_at",
			"received_at",
			"message_status_id",
			"created_at",
			"updated_at",
			"deleted_at",
		})

	now := time.Now() // Assuming 'now' is the current time for the test

	// Mock data row for messages
	MessageSelectQueryRow := testMessagesColumns.AddRow(
		1,                                   // id
		DefaultTestReferenceID,              // reference_id
		DefaultTestConversationID,           // conversation_id
		DefaultFromUserID,                   // from_user_id
		DefaultFromUserID,                   // to_user_id
		DefaultTestMessage,                  // body
		models.NewMessageTypeSMS().ID,       // message_type_id
		now,                                 // sent_at
		now,                                 // received_at
		models.NewMessageStatusPending().ID, // message_status_id
		now,                                 // created_at
		now,                                 // updated_at
		nil,                                 // deleted_at
	)

	return MessageSelectQueryRow
}

// GenerateMockMessageStatus generates a mock message status
// row for testing. It uses a MessageStatusSent
func GenerateMockMessageStatus() *sqlmock.Rows {
	TestAccountStatusesColumns := sqlmock.NewRows(
		[]string{"id", "name"})

	return TestAccountStatusesColumns.AddRow(
		1, models.NewMessageStatusPending().Name,
	)
}

// GenerateMockConversation generates a mock conversation SQL row
// in the 'open' state.
func GenerateMockConversation(open bool) *sqlmock.Rows {
	rows := GenerateMockConversationRowColumns()
	now := time.Now()
	end := &now

	if !open {
		end = nil
	}

	rows.AddRow(1, 1, now, end, now, now, nil)

	return rows
}

// GenerateMockOpenConversations generates a mock open conversation
// SQL response with two rows, with IDs 1 and 2 respectively.
func GenerateMockOpenConversations() *sqlmock.Rows {
	rows := GenerateMockConversationRowColumns()
	now := time.Now()

	rows.AddRow(1, 1, now, nil, now, nil, nil)
	rows.AddRow(2, 1, now, nil, now, nil, nil)
	return rows
}

func GenerateMockConversationRowColumns() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id",
		"user_id",
		"start_time",
		"end_time",
		"created_at",
		"updated_at",
		"deleted_at",
	})
}

func GenerateMockLastAffectedRow() sql.Result {
	return sqlmock.NewResult(1, 1)
}

// GenerateTestConfig creates a Config object with predefined test values
func GenerateTestConfig() config.Config {
	return config.Config{
		OpenAIAPIKey:                 "dummy_openai_api_key",
		DatabaseHost:                 "dummy_database_host",
		DatabaseUser:                 "dummy_database_user",
		DatabasePassword:             "dummy_database_password",
		DatabaseName:                 "dummy_database_name",
		LogLevel:                     100,
		SMSQueueURL:                  "dummy_sms_queue_url",
		TwilioSID:                    "dummy_twilio_sid",
		TwilioAuthToken:              "dummy_twilio_auth_token",
		TwilioPhoneNumber:            "dummy_twilio_phone_number",
		TwilioStatusCallbackURL:      "dummy_twilio_status_callback_url",
		TwilioVerifyServiceSID:       "dummy_twilio_verify_service_sid",
		SNSTopicARN:                  "dummy_sns_topic_arn",
		ChatModelName:                "dummy_chat_model_name",
		ChatModelMaxCompletionTokens: 1000,
		ChatModelFrequencyPenalty:    1.0,
		ChatModelTemperature:         1.0,
	}
}

// SetupMockDB creates and returns a mock SQL DB and a GORM DB.
// This is useful for testing purposes where a real database connection isn't necessary or ideal.
func SetupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	// Retrieve the environment variable to determine the testing environment.
	env := os.Getenv("TEST_ENV")

	if env == "" {
		// If the environment variable isn't set, default to "local".
		env = "mocked"
	}

	// Create a new mock SQL database using the sqlmock package.
	// This mock database simulates the behavior of a real SQL database.
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		// If there's an error in creating the mock database, return the error.
		return nil, nil, err
	}

	var database *gorm.DB

	// Check if the environment is set to use Docker.
	if env == "docker" {
		// In a Docker environment, use the real database configuration.
		database = db.Get(config.Get())
	} else {
		// If not running in Docker, use the sqlmock database.
		dialector := mysql.New(mysql.Config{
			DriverName: "mysql",
			DSN:        "sqlmock_db_0",
			Conn:       mockDB,
		})

		// Set up an expectation for a common initial query made by GORM.
		// This is necessary for GORM to operate correctly with the mock database.
		mock.ExpectQuery("SELECT VERSION()").
			WillReturnRows(sqlmock.NewRows([]string{"version"}).
				AddRow("5.7.31"))

		// Open a GORM database connection with the mock database.
		database, err = gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel((*config.Get()).LogLevel)),
		})
	}

	// Return the GORM database instance, the mock SQL database, and any error that occurred.
	return database, mock, err
}

func SetEnvVars() {
	cfg := GenerateTestConfig()

	_ = os.Setenv("OPENAI_API_KEY", cfg.OpenAIAPIKey)
	_ = os.Setenv("DATABASE_USER", cfg.DatabaseUser)
	_ = os.Setenv("DATABASE_PASSWORD", cfg.DatabasePassword)
	_ = os.Setenv("DATABASE_NAME", cfg.DatabaseName)
	_ = os.Setenv("DATABASE_HOST", cfg.DatabaseHost)
	_ = os.Setenv("TWILIO_AUTH_TOKEN", cfg.TwilioAuthToken)
	_ = os.Setenv("TWILIO_SID", cfg.TwilioSID)
	_ = os.Setenv("TWILIO_PHONE_NUMBER", cfg.TwilioPhoneNumber)
	_ = os.Setenv("LOG_LEVEL", strconv.Itoa(cfg.LogLevel)) // Convert int to string
	_ = os.Setenv("SMS_QUEUE_URL", cfg.SMSQueueURL)
	_ = os.Setenv("TWILIO_STATUS_CALLBACK_URL", cfg.TwilioStatusCallbackURL)
	_ = os.Setenv("TWILIO_VERIFY_SERVICE_SID", cfg.TwilioVerifyServiceSID)
	_ = os.Setenv("SNS_TOPIC_ARN", cfg.SNSTopicARN)
	_ = os.Setenv("CHAT_MODEL_TEMPERATURE", strconv.FormatFloat(float64(cfg.ChatModelTemperature), 'f', -1, 64))
	_ = os.Setenv("CHAT_MODEL_FREQUENCY_PENALTY", strconv.FormatFloat(float64(cfg.ChatModelFrequencyPenalty), 'f', -1, 64))
	_ = os.Setenv("CHAT_MODEL_MAX_COMPLETION_TOKENS", strconv.Itoa(cfg.ChatModelMaxCompletionTokens))
	_ = os.Setenv("CHAT_MODEL_NAME", cfg.ChatModelName)
}
