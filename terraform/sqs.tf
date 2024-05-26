# Topic for fanout of received SMS message
resource "aws_sns_topic" "sms_inbound_topic" {
  name = "sms-inbound-topic"
}

# Queue for outbound sender lambda
resource "aws_sqs_queue" "sms_inbound_queue" {
  name = "sms-inbound-queue"
}

# Queue for the factfinder lambda
resource "aws_sqs_queue" "sms_factfinder_queue" {
  name = "sms-factfinder-queue"
}

# Subscribe Sender Lambda queue to SNS topic
resource "aws_sns_topic_subscription" "sender_subscription" {
  topic_arn = aws_sns_topic.sms_inbound_topic.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.sms_inbound_queue.arn
}

resource "aws_sns_topic_subscription" "factfinder_subscription" {
  topic_arn = aws_sns_topic.sms_inbound_topic.arn
  protocol  = "sqs"
  endpoint  = aws_sqs_queue.sms_factfinder_queue.arn
}

# IAM policy to allow SNS to send messages to SQS
resource "aws_iam_policy" "sns_to_sqs_policy" {
  name = "sns-to-sqs-policy"

  policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Effect : "Allow",
        Action : "sqs:SendMessage",
        Resource : [
          aws_sqs_queue.sms_factfinder_queue.arn,
          aws_sqs_queue.sms_inbound_queue.arn,
        ]
      }
    ]
  })
}

# Permission to allow SQS to invoke sender lambda
resource "aws_lambda_permission" "allow_sender_lambda_sqs" {
  statement_id  = "AllowExecutionFromSQS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.send_sms_lambda.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = aws_sqs_queue.sms_inbound_queue.arn
}

resource "aws_lambda_permission" "allow_factfinder_lambda_sqs" {
  statement_id  = "AllowExecutionFromSQS"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.factfinder_sms_lambda.function_name
  principal     = "sqs.amazonaws.com"
  source_arn    = aws_sqs_queue.sms_factfinder_queue.arn
}

# Invoke the sender when a message is received
resource "aws_lambda_event_source_mapping" "sqs_to_sender_lambda_trigger" {
  event_source_arn  = aws_sqs_queue.sms_inbound_queue.arn
  function_name     = aws_lambda_function.send_sms_lambda.arn
  enabled           = true
}

# Invoke the sender when a message is received
resource "aws_lambda_event_source_mapping" "sqs_to_factfinder_lambda_trigger" {
  event_source_arn = aws_sqs_queue.sms_factfinder_queue.arn
  function_name    = aws_lambda_function.factfinder_sms_lambda.arn
  enabled          = true
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
    aws_security_group.receiver_lambda_sg.id,
    aws_security_group.factfinder_lambda_sg.id
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
          "sqs:ChangeMessageVisibility"
        ],
        Resource : [
          aws_sqs_queue.sms_inbound_queue.arn,
          aws_sqs_queue.sms_factfinder_queue.arn
        ],
        Effect: "Allow",
      },
    ],
  })
}


# Policy to allow SNS to send messages to sms_inbound_queue
resource "aws_sqs_queue_policy" "sms_inbound_queue_policy" {
  queue_url = aws_sqs_queue.sms_inbound_queue.id

  policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Effect : "Allow",
        Principal : "*",
        Action : [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
          "sqs:ChangeMessageVisibility"
        ],
        Resource : aws_sqs_queue.sms_inbound_queue.arn,
        Condition : {
          ArnEquals : {
            "aws:SourceArn" : aws_sns_topic.sms_inbound_topic.arn
          }
        }
      }
    ]
  })
}

resource "aws_sqs_queue_policy" "sms_factfinder_queue_policy" {
  queue_url = aws_sqs_queue.sms_factfinder_queue.id

  policy = jsonencode({
    Version : "2012-10-17",
    Statement : [
      {
        Effect : "Allow",
        Principal : "*",
        Action : [
          "sqs:SendMessage",
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes",
          "sqs:ChangeMessageVisibility"

        ],
        Resource : aws_sqs_queue.sms_factfinder_queue.arn,
        Condition : {
          ArnEquals : {
            "aws:SourceArn" : aws_sns_topic.sms_inbound_topic.arn
          }
        }
      }
    ]
  })
}
