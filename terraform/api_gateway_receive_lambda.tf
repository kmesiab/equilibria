#
# Sets up the URL path for /{env}/sms-receive
#
resource "aws_api_gateway_resource" "api_route_sms_receiver" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  parent_id   = aws_api_gateway_rest_api.api_gateway.root_resource_id
  path_part   = "sms-receive"

  lifecycle {
    create_before_destroy = true
  }
}

#
# POST /sms-receive
#
resource "aws_api_gateway_method" "sms_receive_post_method" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  resource_id   = aws_api_gateway_resource.api_route_sms_receiver.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "sms_receive_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api_gateway.id
  resource_id             = aws_api_gateway_resource.api_route_sms_receiver.id
  http_method             = "POST"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.receiver_sms_lambda.invoke_arn
}
