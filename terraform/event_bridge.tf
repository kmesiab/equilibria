resource "aws_cloudwatch_event_rule" "nudger_event_rule" {
  name                = "nudger-event-rule"
  description         = "Triggers nudger Lambda function"
  schedule_expression = "rate(1 hour)"
}

resource "aws_cloudwatch_event_target" "nudger_event_target" {
  rule = aws_cloudwatch_event_rule.nudger_event_rule.name
  arn  = aws_lambda_function.nudger_lambda.arn
}

resource "aws_lambda_permission" "allow_event_bridge" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.nudger_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.nudger_event_rule.arn
}
