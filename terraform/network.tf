resource "aws_ecs_cluster" "devbot-cluster" {
  name = "devbot-cluster"
}

data "aws_availability_zones" "aws-az" {
  state = "available"
}

# Internet VPC
resource "aws_vpc" "devbot" {
  cidr_block           = local.vpc_cidr
  instance_tenancy     = "default"
  enable_dns_support   = "true"
  enable_dns_hostnames = "true"

  enable_classiclink               = false


  assign_generated_ipv6_cidr_block = "false"

  tags = merge(
  local.tags,
  {
    Name = local.vpc_name
  }
  )
}

resource "aws_subnet" "private" {
  count = length(data.aws_availability_zones.aws-az.names)
  vpc_id = aws_vpc.devbot.id
  cidr_block = cidrsubnet(aws_vpc.devbot.cidr_block, 8, count.index + 1)
  availability_zone = data.aws_availability_zones.aws-az.names[count.index]
  map_public_ip_on_launch = true
  tags = {
    Name = "${local.application}-subnet-${count.index + 1}"
  }
}

# create internet gateway
resource "aws_internet_gateway" "aws-igw" {
  vpc_id = aws_vpc.devbot.id
  tags = {
    Name = "${local.application}-igw"
  }
}
# create routes
resource "aws_route_table" "aws-route-table" {
  vpc_id = aws_vpc.devbot.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.aws-igw.id
  }
  tags = {
    Name = "${local.application}-route-table"
  }
}
resource "aws_main_route_table_association" "aws-route-table-association" {
  vpc_id = aws_vpc.devbot.id
  route_table_id = aws_route_table.aws-route-table.id
}