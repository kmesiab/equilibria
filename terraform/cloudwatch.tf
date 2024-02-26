resource "aws_cloudwatch_log_group" "api_gateway_logs" {
  name = "/equilibria/dev/equilibria-gateway"

  retention_in_days = 1
}

# resource "aws_cloudwatch_log_group" "mange_user_lambda_log_group" {
#  name = "/aws/lambda/manageUserFunction"

#   retention_in_days = 1
# }

# resource "aws_cloudwatch_log_group" "sms_receive_lambda_log_group" {
#   name = "/aws/lambda/smsReceiveFunction"

#   retention_in_days = 1
# }


# resource "aws_cloudwatch_log_group" "sms_send_lambda_log_group" {
#   name = "/aws/lambda/smsSendFunction"

#   retention_in_days = 1
# }


# resource "aws_cloudwatch_log_group" "sms_status_lambda_log_group" {
#   name = "/aws/lambda/smsStatusFunction"

#  retention_in_days = 1
# }
