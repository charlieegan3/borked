resource "aws_dynamodb_table" "jobs" {
  name           = "${var.project}-jobs"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "JobId"

  attribute {
    name = "JobId"
    type = "S"
  }
}

resource "aws_dynamodb_table" "counts" {
  name           = "${var.project}-counts"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "Count"

  attribute {
    name = "Count"
    type = "S"
  }
}
