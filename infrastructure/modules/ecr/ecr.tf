variable "services" {
  default = [
    "user-service",
    "product-service",
    "purchase-service",
    "file-service"
  ]
}

resource "aws_ecr_repository" "microservices" {
  for_each = toset(var.services)
  name     = each.value
}