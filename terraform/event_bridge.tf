resource "aws_cloudwatch_event_rule" "robo_trade1_lambda_event_rule" {
  name                = "robo_trade1-lambda-event-rule"
  description         = "指定時間毎に実行"
  schedule_expression = "rate(10 minutes)"
}
 
resource "aws_cloudwatch_event_target" "robo_trade1_lambda_target" {
  arn  = module.robo_trade1_lambda.lambda_function_arn
  rule = aws_cloudwatch_event_rule.robo_trade1_lambda_event_rule.name
}
 
resource "aws_lambda_permission" "robo_trade1_lambda_permission" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = module.robo_trade1_lambda.lambda_function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.robo_trade1_lambda_event_rule.arn
}