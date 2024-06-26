locals {
  lambda_environment_variables = {
    DATABASE_HOST                    = aws_db_instance.mysql.address
    DATABASE_USER                    = aws_db_instance.mysql.username
    DATABASE_PASSWORD                = aws_db_instance.mysql.password
    DATABASE_NAME                    = aws_db_instance.mysql.db_name
    SMS_QUEUE_URL                    = aws_sqs_queue.sms_inbound_queue.url
    OPENAI_API_KEY                   = aws_ssm_parameter.openai_api_key.value
    TWILIO_SID                       = aws_ssm_parameter.twilio_sid.value
    TWILIO_AUTH_TOKEN                = aws_ssm_parameter.twilio_auth_token.value
    TWILIO_PHONE_NUMBER              = aws_ssm_parameter.twilio_phone_number.value
    TWILIO_STATUS_CALLBACK_URL       = aws_ssm_parameter.twilio_status_callback_url.value
    TWILIO_VERIFY_SERVICE_SID        = aws_ssm_parameter.twilio_verify_service_sid.value
    CHAT_MODEL_NAME                  = aws_ssm_parameter.chat_model_name.value
    CHAT_MODEL_TEMPERATURE           = aws_ssm_parameter.chat_model_temperature.value
    CHAT_MODEL_MAX_COMPLETION_TOKENS = aws_ssm_parameter.chat_model_max_completion_tokens.value
    CHAT_MODEL_FREQUENCY_PENALTY     = aws_ssm_parameter.chat_model_frequency_penalty.value
    SNS_TOPIC_ARN                    = aws_sns_topic.sms_inbound_topic.arn
  }
}
