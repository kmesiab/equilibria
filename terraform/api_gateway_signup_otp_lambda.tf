#
# Sets up the URL path for /{env}/signup-otp
#
resource "aws_api_gateway_resource" "api_route_signup_otp" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  parent_id   = aws_api_gateway_rest_api.api_gateway.root_resource_id
  path_part   = "signup-otp"

  lifecycle {
    create_before_destroy = true
  }
}
#
# POST /signup-otp
#
resource "aws_api_gateway_method" "signup_otp_post_method" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  resource_id   = aws_api_gateway_resource.api_route_signup_otp.id
  http_method   = "POST"
  authorization = "NONE"
}
#
# PUT /signup-otp
#
resource "aws_api_gateway_method" "signup_otp_put_method" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  resource_id   = aws_api_gateway_resource.api_route_signup_otp.id
  http_method   = "PUT"
  authorization = "NONE"
}
#
# OPTIONS /signup-otp
#
resource "aws_api_gateway_method" "signup_otp_options_method" {
  rest_api_id   = aws_api_gateway_rest_api.api_gateway.id
  resource_id   = aws_api_gateway_resource.api_route_signup_otp.id
  http_method   = "OPTIONS"
  authorization = "NONE"
}
#
# OPTIONS /signup-otp
#
resource "aws_api_gateway_method_response" "signup_otp_options_method_response" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  resource_id = aws_api_gateway_resource.api_route_signup_otp.id
  http_method = aws_api_gateway_method.signup_otp_options_method.http_method
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin"  = true
  }
}

resource "aws_api_gateway_integration_response" "signup_otp_options_integration_response" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  resource_id = aws_api_gateway_resource.api_route_signup_otp.id
  http_method = aws_api_gateway_method.signup_otp_options_method.http_method
  status_code = aws_api_gateway_method_response.signup_otp_options_method_response.status_code

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token,X-Amz-User-Agent'"
    "method.response.header.Access-Control-Allow-Methods" = "'DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT'"
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
  }

  response_templates = {
    "application/json" = ""
  }
}


resource "aws_api_gateway_integration" "signup_otp_post_lambda_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api_gateway.id
  resource_id             = aws_api_gateway_resource.api_route_signup_otp.id
  http_method             = "POST"
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.signup_otp_lambda.invoke_arn
}

resource "aws_api_gateway_integration" "signup_otp_put_integration" {
  rest_api_id             = aws_api_gateway_rest_api.api_gateway.id
  resource_id             = aws_api_gateway_resource.api_route_signup_otp.id
  http_method             = aws_api_gateway_method.signup_otp_put_method.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.signup_otp_lambda.invoke_arn
}

resource "aws_api_gateway_integration" "signup_otp_options_integration" {
  rest_api_id = aws_api_gateway_rest_api.api_gateway.id
  resource_id = aws_api_gateway_resource.api_route_signup_otp.id
  http_method = aws_api_gateway_method.signup_otp_options_method.http_method
  type        = "MOCK"

  request_templates = {
    "application/json" = "{\"statusCode\": 200}"
  }
}
