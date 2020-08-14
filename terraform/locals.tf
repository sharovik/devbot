locals {
  account_id = data.aws_caller_identity.current.account_id
  application = "devbot-application"
  env = "development"
  tags = {
    Environment = terraform.workspace
    Project     = "shared"
    Terraform   = "base"
  }
  vpc_name = "devbot-${terraform.workspace}"
  vpc_id       = aws_vpc.devbot.id
  vpc_cidr     = "10.0.0.0/16"
  namespace = replace(
    title("${local.application}-${terraform.workspace}"),
    "-",
    "",
  )
  subnet_private_availability_zone_map = {
    "0" = "us-east-1b"
    "1" = "us-east-1d"
    #f does not support SES
    "2" = "us-east-1f"
  }
}