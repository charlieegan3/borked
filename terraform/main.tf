variable project {
  default = "borked"
}

variable domain {
  default = "borked.charlieegan3.com"
}

data "aws_acm_certificate" "default" {
  domain   = "charlieegan3.com"
  statuses = ["ISSUED"]
}

data "aws_caller_identity" "current" {}

variable "region" {
  default = "us-east-1"
}

terraform {
  backend "s3" {
    bucket = "charlieegan3-www-terraform-state"
    region = "us-east-1"
    key    = "borked.tfstate"
  }
}
