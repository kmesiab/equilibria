#
# Sets up the URL path for /{env}/sms-status
#
resource "aws_api_gateway_resource" "api_route_sms_status" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  parent_id   = aws_api_gateway_rest_api.api_gateway.root_resource_id
  path_part   = "sms-status"
}

#
# POST /sms-status
#
resource "aws_api_gateway_method" "sms_status_post_method" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id             # The API Gateway
  resource_id   = aws_api_gateway_resource.api_route_sms_status.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "sms_status_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api_gateway.id
  resource_id             = aws_api_gateway_resource.api_route_sms_status.id
  http_method             = "POST"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.sms_status_lambda.invoke_arn
}
