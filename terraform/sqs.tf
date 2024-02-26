resource "aws_sqs_queue" "sms_inbound_queue" {
  name = "sms-inbound-queue"
}

resource "aws_lambda_permission" "allow_sender_lambda_sqs" {
  statement_id  = "AllowExecutionFromSQS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.send_sms_lambda.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = aws_sqs_queue.sms_inbound_queue.arn
}

# Invoke the sender when a message is received
resource "aws_lambda_event_source_mapping" "sms_responder" {
  event_source_arn  = aws_sqs_queue.sms_inbound_queue.arn
  function_name     = aws_lambda_function.send_sms_lambda.arn
  enabled           = true
  batch_size        = 10
}

# This endpoint lets us interact with the queue without having to expose it publicly
resource "aws_vpc_endpoint" "sqs_endpoint" {
  vpc_id            = aws_vpc.my_vpc.id
  service_name      = "com.amazonaws.${var.region}.sqs"
  vpc_endpoint_type = "Interface"

  subnet_ids = [aws_subnet.receiver_subnet.id,aws_subnet.outbound_subnet.id]

  private_dns_enabled = true
  security_group_ids  = [
    aws_security_group.sender_lambda_sg.id,
    aws_security_group.receiver_lambda_sg.id
  ]
}

# Policy that allows interaction with the lambda
resource "aws_iam_policy" "lambda_sqs_policy" {
  name        = "lambda-sqs-policy"
  description = "IAM policy for allowing Lambda to interact with SQS queue"

  policy = jsonencode({
    Version: "2012-10-17",
    Statement: [
      {
        Action: [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
        ],
        Resource: aws_sqs_queue.sms_inbound_queue.arn,
        Effect: "Allow",
      },
    ],
  })
}

