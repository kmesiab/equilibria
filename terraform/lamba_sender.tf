resource "aws_lambda_function" "send_sms_lambda" {
  function_name = "smsSendFunction"
  runtime       = "go1.x"
  handler       = "main"
  timeout       = 30
  filename      = "../build/send_sms.zip"
  role          = aws_iam_role.lambda_execution_role.arn

  environment {
    variables = local.lambda_environment_variables
  }

}

resource "aws_security_group" "sender_lambda_sg" {
  name        = "sender_lambda_sg"
  description = "Security group for Lambda function"
  vpc_id      = aws_vpc.my_vpc.id

  # Outbound rule to allow Lambda to communicate with the RDS instance
  egress {
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.my_vpc.cidr_block]  # VPC CIDR block
  }

  # Outbound rule to allow Lambda to get responses from the RDS instance
  ingress {
    from_port   = 3306
    to_port     = 3306
    protocol    = "tcp"
    cidr_blocks = [aws_vpc.my_vpc.cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [
      aws_vpc.my_vpc.cidr_block,
      "0.0.0.0/0"
    ]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [
      aws_vpc.my_vpc.cidr_block,
      "0.0.0.0/0"
    ]
  }

  tags = {
    Name        = "sender_lambda_sg"
    Description = "Security group for lambda functions requiring outbound internet access"
  }
}
