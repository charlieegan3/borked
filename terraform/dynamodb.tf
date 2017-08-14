resource "aws_dynamodb_table" "jobs" {
  name           = "BorkedJobs"
  read_capacity  = 20
  write_capacity = 20
  hash_key       = "JobId"

  attribute {
    name = "JobId"
    type = "S"
  }
}
