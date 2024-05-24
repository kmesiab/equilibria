resource "aws_ssm_parameter" "openai_api_key" {
  name  = "/config/OPENAI_API_KEY"
  type  = "SecureString"
  value = var.open_ai_api_key
}

resource "aws_ssm_parameter" "mysql_root_password" {
  name  = "/config/MYSQL_ROOT_PASSWORD"
  type  = "SecureString"
  value = var.mysql_root_password
}

resource "aws_ssm_parameter" "mysql_root_user" {
  name  = "/config/MYSQL_ROOT_USER"
  type  = "SecureString"
  value = var.mysql_root_user
}

resource "aws_ssm_parameter" "database_user" {
  name  = "/config/DATABASE_USER"
  type  = "SecureString"
  value = var.database_user
}

resource "aws_ssm_parameter" "database_password" {
  name  = "/config/DATABASE_PASSWORD"
  type  = "SecureString"
  value = var.database_password
}

resource "aws_ssm_parameter" "database_name" {
  name  = "/config/DATABASE_NAME"
  type  = "SecureString"
  value = var.database_name
}

resource "aws_ssm_parameter" "database_host" {
  name  = "/config/DATABASE_HOST"
  type  = "SecureString"
  value = aws_db_instance.mysql.address
}

resource "aws_ssm_parameter" "sms_inbound_queue_url" {
  name  = "/config/SMS_QUEUE_URL"
  type  = "SecureString"
  value = aws_sqs_queue.sms_inbound_queue.url
}

resource "aws_ssm_parameter" "twilio_sid" {
  name  = "/config/TWILIO_SID"
  type  = "SecureString"
  value = var.twilio_sid
}

resource "aws_ssm_parameter" "twilio_auth_token" {
  name  = "/config/TWILIO_AUTH_TOKEN"
  type  = "SecureString"
  value = var.twilio_auth_token
}

resource "aws_ssm_parameter" "twilio_phone_number" {
  name  = "/config/TWILIO_PHONE_NUMBER"
  type  = "SecureString"
  value = var.twilio_phone_number
}

resource "aws_ssm_parameter" "twilio_status_callback_url" {
  name  = "/config/TWILIO_STATUS_CALLBACK_URL"
  type  = "SecureString"
  value = var.twilio_status_callback_url
}

resource "aws_ssm_parameter" "twilio_verify_service_sid" {
  name  = "/config/TWILIO_VERIFY_SERVICE_SID"
  type  = "SecureString"
  value = var.twilio_verify_service_sid
}

resource "aws_ssm_parameter" "api_endpoint" {
  name  = "/config/API_HOSTNAME"
  type  = "SecureString"
  value = "https://${aws_api_gateway_rest_api.api_gateway.id}.execute-api.${var.region}.amazonaws.com/${aws_api_gateway_stage.api_stage.stage_name}"
}

resource "aws_ssm_parameter" "chat_model_name" {
  name  = "/config/CHAT_MODEL_NAME"
  type  = "SecureString"
  value = var.chat_model_name
}

resource "aws_ssm_parameter" "chat_model_temperature" {
  name  = "/config/CHAT_MODEL_TEMPERATURE"
  type  = "SecureString"
  value = var.chat_model_temperature
}

resource "aws_ssm_parameter" "chat_model_max_completion_tokens" {
  name  = "/config/CHAT_MODEL_MAX_COMPLETION_TOKENS"
  type  = "SecureString"
  value = var.chat_model_max_completion_tokens
}

resource "aws_ssm_parameter" "chat_model_frequency_penalty" {
  name  = "/config/CHAT_MODEL_FREQUENCY_PENALTY"
  type  = "SecureString"
  value = var.chat_model_frequency_penalty
}
