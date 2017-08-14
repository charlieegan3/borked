resource "aws_lambda_permission" "apigw_lambda_count" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.count.arn}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.count.id}/*/*/*"
}

resource "aws_lambda_function" "count" {
  filename         = "../count_function/handler.zip"
  function_name    = "${var.project}-count"
  role             = "${aws_iam_role.lambda.arn}"
  handler          = "handler.Handle"
  runtime          = "python2.7"
  timeout          = "30"
  source_code_hash = "${base64sha256(file("../count_function/handler.zip"))}"
  memory_size      = "512"
}
