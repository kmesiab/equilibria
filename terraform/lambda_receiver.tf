# Lambda Function for SMS Status
resource "aws_lambda_function" "receiver_sms_lambda" {
  function_name = "smsReceiveFunction"
  runtime       = "go1.x"
  handler       = "main"
  timeout       = 30
  filename      = "../build/receive_sms.zip"
  role          = aws_iam_role.lambda_execution_role.arn

  environment {
    variables = local.lambda_environment_variables
  }

  vpc_config {
    subnet_ids         = [aws_subnet.receiver_subnet.id, aws_subnet.outbound_subnet.id]
    security_group_ids = [aws_security_group.receiver_lambda_sg.id]
  }
}


resource "aws_security_group" "receiver_lambda_sg" {
  name        = "sg_receiver_subnet"
  description = "Security group for Lambda in receiver subject"
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

  # Include an egress to reach the SQS queue
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [aws_vpc.my_vpc.cidr_block]
  }

  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = [aws_vpc.my_vpc.cidr_block]
  }
}
