resource "aws_api_gateway_usage_plan" "MyUsagePlan" {
  name        = "throttling-plan"
  description = "Throttle requests into the borked API"

  api_stages {
    api_id = "${aws_api_gateway_rest_api.default.id}"
    stage  = "${var.api_gateway_stage}"
  }

  throttle_settings {
    burst_limit = 10
    rate_limit  = 25
  }
}
