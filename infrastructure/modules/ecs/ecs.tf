resource "aws_ecs_cluster" "tutuplapak_cluster" {
  name = "tutuplapak_cluster"
}

resource "aws_default_vpc" "default_vpc" {}

resource "aws_subnet" "private_subnet" {
  availability_zone = "ap-southeast-1a"
}

resource "aws_iam_role" "ecs_task_execution_role" {
  name               = "tutuplapak-ecs-task-execution-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role_policy.json
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role_policy" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

resource "aws_alb" "application_load_balancer" {
  name               = "tutuplapak-alb"
  load_balancer_type = "application"
  subnets = [
    "${aws_subnet.private_subnet.id}"
  ]
  security_groups = ["${aws_security_group.load_balancer_security_group.id}"]
}

resource "aws_lb_target_group" "target_group" {
  name        = "tutuplapak-tg"
  port        = "80"
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_default_vpc.default_vpc.id
}

resource "aws_lb_listener" "http_listener" {
  load_balancer_arn = aws_alb.application_load_balancer.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.target_group.arn
  }
}

resource "aws_security_group" "load_balancer_security_group" {
  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "ecs_services" {
  # Only allow traffic from ALB
  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = [aws_security_group.load_balancer_security_group.id]
  }

  # Allow all outbound traffic (e.g., for ECS tasks to call APIs or download updates)
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Service Start

variable "services" {
  default = ["user-service", "product-service", "purchase-service", "file-service"]
}

variable "service_configs" {
  default = {
    "service-user"     = { container_port = 8081, cpu = 1024, memory = 4096, path_pattern = "/v1/user/*" }
    "service-product"  = { container_port = 8082, cpu = 2048, memory = 8192, path_pattern = "/v1/product/*" }
    "service-purchase" = { container_port = 8083, cpu = 1024, memory = 4096, path_pattern = "/v1/purchase/*" }
    "service-file"     = { container_port = 8084, cpu = 1024, memory = 4096, path_pattern = "/v1/file/*" }
  }
}

resource "aws_ecs_task_definition" "task_definitions" {
  for_each = var.services

  family                   = each.value
  container_definitions    = jsonencode([
    {
      name  = each.value
      image = "${aws_account_id}.dkr.ecr.ap-southeast-1.amazonaws.com/${each.value}:latest"
      memory = var.service_configs[each.value].memory
      cpu    = var.service_configs[each.value].cpu
      essential = true
      portMappings = [
        {
          containerPort = var.service_configs[each.value].container_port
          hostPort      = var.service_configs[each.value].container_port
        }
      ]
    }
  ])
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  memory                   = var.service_configs[each.value].memory
  cpu                      = var.service_configs[each.value].cpu
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
}

resource "aws_ecs_service" "ecs_services" {
  for_each = toset(var.services)

  name            = each.value
  cluster         = aws_ecs_cluster.tutuplapak_cluster.id
  task_definition = aws_ecs_task_definition.task_definitions[each.value].arn
  launch_type     = "FARGATE"
  desired_count   = 1

  load_balancer {
    target_group_arn = aws_lb_target_group.target_groups[each.key].arn
    container_name   = aws_ecs_task_definition.task_definitions[each.key].family
    container_port   = var.service_configs[each.key].container_port
  }

  network_configuration {
    subnets          = [aws_subnet.private_subnet.id]
    assign_public_ip = false
    security_groups  = [aws_security_group.ecs_services.id]
  }
}

resource "aws_lb_target_group" "target_groups" {
  for_each = toset(var.services)

  name        = "${each.value}-tg"
  port        = var.service_configs[each.value].container_port
  protocol    = "HTTP"
  target_type = "ip"
}

resource "aws_lb_listener_rule" "path_based_routing" {
  for_each = toset(var.services)

  listener_arn = aws_lb_listener.http_listener.arn
  priority     = 10 + index(var.services, each.key)

  conditions {
    field  = "path-pattern"
    values = [var.service_configs[each.key].path_pattern]
  }

  actions {
    type             = "forward"
    target_group_arn = aws_lb_target_group.target_groups[each.value].arn
  }
}

# Service End

