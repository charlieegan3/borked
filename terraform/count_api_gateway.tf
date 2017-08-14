output "count_api_endpoint" {
  value = "${aws_api_gateway_deployment.count.invoke_url}/${var.count_api_gateway_path_part}"
}

variable "count_api_gateway_stage" {
  default = "production"
}

variable "count_api_gateway_path_part" {
  default = "count"
}

resource "aws_api_gateway_rest_api" "count" {
  name        = "${var.project}"
  description = "API for the ${var.project} lambda function (count)"
}

resource "aws_api_gateway_deployment" "count" {
  depends_on  = ["aws_api_gateway_integration.count"]
  rest_api_id = "${aws_api_gateway_rest_api.count.id}"
  stage_name  = "${var.count_api_gateway_stage}"

  // https://github.com/hashicorp/terraform/issues/6613
  stage_description = "${md5(file("count_api_gateway.tf"))}"
}

resource "aws_api_gateway_method_settings" "count" {
  rest_api_id = "${aws_api_gateway_rest_api.count.id}"
  stage_name  = "${var.count_api_gateway_stage}"
  method_path = "${aws_api_gateway_resource.count.path_part}/*"

  settings {
    metrics_enabled = true
    logging_level   = "INFO"
  }

  depends_on = ["aws_api_gateway_deployment.count"]
}

resource "aws_api_gateway_resource" "count" {
  rest_api_id = "${aws_api_gateway_rest_api.count.id}"
  parent_id   = "${aws_api_gateway_rest_api.count.root_resource_id}"
  path_part   = "${var.count_api_gateway_path_part}"
}

resource "aws_api_gateway_method" "count" {
  rest_api_id   = "${aws_api_gateway_rest_api.count.id}"
  resource_id   = "${aws_api_gateway_resource.count.id}"
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "count" {
  rest_api_id             = "${aws_api_gateway_rest_api.count.id}"
  resource_id             = "${aws_api_gateway_resource.count.id}"
  http_method             = "${aws_api_gateway_method.count.http_method}"
  type                    = "AWS_PROXY"
  uri                     = "arn:aws:apigateway:${var.region}:lambda:path/2015-03-31/functions/${aws_lambda_function.count.arn}/invocations"
  integration_http_method = "POST"
}

resource "aws_api_gateway_method_response" "count_response_method" {
  rest_api_id = "${aws_api_gateway_rest_api.count.id}"
  resource_id = "${aws_api_gateway_resource.count.id}"
  http_method = "${aws_api_gateway_integration.count.http_method}"
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Origin" = true
  }
}

resource "aws_api_gateway_integration_response" "count_response_method_integration" {
  rest_api_id = "${aws_api_gateway_rest_api.count.id}"
  resource_id = "${aws_api_gateway_resource.count.id}"
  http_method = "${aws_api_gateway_method_response.count_response_method.http_method}"
  status_code = "${aws_api_gateway_method_response.count_response_method.status_code}"
}

resource "aws_api_gateway_method" "count_resource_options" {
  rest_api_id   = "${aws_api_gateway_rest_api.count.id}"
  resource_id   = "${aws_api_gateway_resource.count.id}"
  http_method   = "OPTIONS"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "count_resource_options" {
  rest_api_id = "${aws_api_gateway_rest_api.count.id}"
  resource_id = "${aws_api_gateway_resource.count.id}"
  http_method = "${aws_api_gateway_method.count_resource_options.http_method}"
  type        = "MOCK"

  request_templates = {
    "application/json" = <<PARAMS
      { "statusCode": 200 }
    PARAMS
  }
}

resource "aws_api_gateway_integration_response" "count_resource_options" {
  rest_api_id = "${aws_api_gateway_rest_api.count.id}"
  resource_id = "${aws_api_gateway_resource.count.id}"
  http_method = "${aws_api_gateway_method.count_resource_options.http_method}"
  status_code = "200"

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods" = "'POST,OPTIONS,GET,PUT,PATCH,DELETE'"
    "method.response.header.Access-Control-Allow-Origin"  = "'https://${var.domain}'"
  }

  depends_on = ["aws_api_gateway_integration.count_resource_options"]
}

resource "aws_api_gateway_method_response" "count_resource_options" {
  rest_api_id = "${aws_api_gateway_rest_api.count.id}"
  resource_id = "${aws_api_gateway_resource.count.id}"
  http_method = "OPTIONS"
  status_code = "200"

  response_models = {
    "application/json" = "Empty"
  }

  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = true
    "method.response.header.Access-Control-Allow-Methods" = true
    "method.response.header.Access-Control-Allow-Origin"  = true
  }
}
