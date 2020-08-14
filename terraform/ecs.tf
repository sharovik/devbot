provider "aws" {
  region                  = "us-east-1"
  shared_credentials_file = "~/.aws/credentials"
  profile                 = "free-tier"
}

resource "aws_ecs_task_definition" "main" {
  family                   = "${local.application}-${terraform.workspace}"
  container_definitions    = jsonencode([local.container_definition])
  requires_compatibilities = ["FARGATE"]
  cpu    = 256
  memory = 512

  task_role_arn      = aws_iam_role.task.arn
  execution_role_arn = "arn:aws:iam::${local.account_id}:role/ecsTaskExecutionRole"
  network_mode       = "awsvpc"
}

resource "aws_ecs_service" "main" {
  name = local.application

  cluster = aws_ecs_cluster.devbot-cluster.arn

  task_definition                    = "${aws_ecs_task_definition.main.family}:${aws_ecs_task_definition.main.revision}"
  desired_count                      = terraform.workspace == "production" ? 1 : 1
  deployment_minimum_healthy_percent = 50
  launch_type                        = "FARGATE"

  network_configuration {
    subnets = aws_subnet.private.*.id
    assign_public_ip = true
  }
}

resource "aws_cloudwatch_log_group" "main" {
  name              = "/ecs/${local.application}-${terraform.workspace}"
  retention_in_days = 30

  tags = local.tags
}