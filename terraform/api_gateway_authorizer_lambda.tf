resource "aws_api_gateway_authorizer" "authorizer" {
  name                   = "authorizerFunction"
  rest_api_id            = aws_api_gateway_rest_api.api_gateway.id
  authorizer_uri         = aws_lambda_function.authorizer_lambda.invoke_arn
  authorizer_credentials = aws_iam_role.lambda_execution_role.arn

  identity_source = "method.request.header.Authorization"

  type = "TOKEN"
}
