locals {
  lambda_environment_variables = {
    DATABASE_HOST       = aws_db_instance.mysql.address
    DATABASE_USER       = aws_db_instance.mysql.username
    DATABASE_PASSWORD   = aws_db_instance.mysql.password
    DATABASE_NAME       = aws_db_instance.mysql.db_name
    SMS_QUEUE_URL       = aws_sqs_queue.sms_inbound_queue.url
    OPENAI_API_KEY      = aws_ssm_parameter.openai_api_key.value
    TWILIO_SID          = aws_ssm_parameter.twilio_sid.value
    TWILIO_AUTH_TOKEN   = aws_ssm_parameter.twilio_auth_token.value
    TWILIO_PHONE_NUMBER = aws_ssm_parameter.twilio_phone_number.value
    TWILIO_STATUS_CALLBACK_URL = aws_ssm_parameter.twilio_status_callback_url.value
    TWILIO_VERIFY_SERVICE_SID = aws_ssm_parameter.twilio_verify_service_sid.value
  }
}
