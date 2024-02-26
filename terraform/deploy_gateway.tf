#
# The gateway deployment stage for prod
#
resource "aws_api_gateway_deployment" "api_deployment" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  stage_name  = "nonprod"
  depends_on  = [
    aws_api_gateway_integration.sms_receive_lambda_integration,
    aws_api_gateway_integration.sms_status_lambda_integration,
    aws_api_gateway_integration.login_lambda_integration,
    aws_api_gateway_integration.signup_otp_post_lambda_integration,
    aws_api_gateway_integration.signup_otp_put_integration,
    aws_api_gateway_integration.signup_otp_options_integration,
    aws_api_gateway_integration.manage_user_get_integration,
    aws_api_gateway_integration.manage_user_post_lambda_integration,
    aws_api_gateway_integration.manage_user_put_integration,
    aws_api_gateway_integration.login_options_integration,
    aws_api_gateway_integration.manage_user_options_integration,
  ]

  triggers = {
    redeployment = timestamp()
  }

  lifecycle {
    create_before_destroy = true
  }
}

#
# API Gateway Stage Configuration for dev
#
resource "aws_api_gateway_stage" "api_stage" {
  stage_name           = "dev"
  rest_api_id          = aws_api_gateway_rest_api.api_gateway.id
  deployment_id        = aws_api_gateway_deployment.api_deployment.id
  xray_tracing_enabled = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway_logs.arn
    format          = "{\"requestId\":\"$context.requestId\", \"ip\":\"$context.identity.sourceIp\", \"caller\":\"$context.identity.caller\", \"user\":\"$context.identity.user\", \"requestTime\":\"$context.requestTime\", \"httpMethod\":\"$context.httpMethod\", \"resourcePath\":\"$context.resourcePath\", \"status\":\"$context.status\", \"protocol\":\"$context.protocol\", \"responseLength\":\"$context.responseLength\"}"
  }
}
