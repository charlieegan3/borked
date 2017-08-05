resource "aws_lambda_permission" "apigw_lambda" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = "${aws_lambda_function.lambda.arn}"
  principal     = "apigateway.amazonaws.com"

  source_arn = "arn:aws:execute-api:${var.region}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.test.id}/*/*/*"
}

resource "aws_lambda_function" "lambda" {
  filename         = "../handler.zip"
  function_name    = "mylambda"
  role             = "${aws_iam_role.role.arn}"
  handler          = "handler.Handle"
  runtime          = "python2.7"
  timeout          = "30"
  source_code_hash = "${base64sha256(file("../handler.zip"))}"
}

resource "aws_iam_role" "role" {
  name = "myrole"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy" "test_policy" {
  name = "test_policy"
  role = "${aws_iam_role.role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "cloudwatch:*",
        "logs:*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}
