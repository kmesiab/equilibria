output "db_endpoint" {
  description = "The connection endpoint for the RDS instance"
  value       = aws_db_instance.mysql.endpoint
}

# Output URL
output "api-url" {
  value       = "https://${aws_api_gateway_rest_api.api_gateway.id}.execute-api.${var.region}.amazonaws.com/${aws_api_gateway_stage.api_stage.stage_name}"
  description = "The URL to invoke the Lambda function via API Gateway"
}

output "sms_inbound_queue_url" {
  value = aws_sqs_queue.sms_inbound_queue.url
}

output "sns_topic_arn" {
  value = aws_sns_topic.sms_inbound_topic.arn
}
