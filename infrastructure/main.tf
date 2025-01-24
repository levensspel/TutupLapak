terraform {
  required_version = "~> 1.10.5"

  backend "s3" {
    bucket         = "infra"
    key            = "tf/terraform.tfstate"
    region         = "ap-southeast-1"
    dynamodb_table = "tfstate-tutuplapak"
    encrypt        = true
  }

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.84.0"
    }
  }
}

module "tf-state" {
  source = "./modules/tf-state"
}

module "ecrRepo" {
  source = "./modules/ecr"
}

module "ecsCluster" {
  source = "./modules/ecs"
}