resource "aws_iam_role" "task" {
  name               = "${local.namespace}Role"
  description        = "Allows ECS tasks to call AWS services on your behalf."
  assume_role_policy = file("policies/ecs_role.json")

  tags = local.tags
}

data "template_file" "devbot_ecs" {
  template = file("policies/devbot_ecs.json")

  vars = {
    account_id       = local.account_id
    env_suffix       = "-${terraform.workspace}"
    env_suffix_fixed = "-development"
    cluster_arn      = aws_ecs_cluster.devbot-cluster.arn
  }
}

resource "aws_iam_policy" "devbot_ecs" {
  name   = "DevbotDevelopmentECS"
  policy = data.template_file.devbot_ecs.rendered
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name = "ecsTaskExecutionRole"

  assume_role_policy = file("policies/role.json")
}

resource "aws_iam_role_policy_attachment" "ecs-task-execution-role-policy-attachment" {
  role = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_iam_role_policy_attachment" "task_s3" {
  role = aws_iam_role.task.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonS3FullAccess"
}