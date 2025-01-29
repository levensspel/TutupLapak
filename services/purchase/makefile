# Required environment variables check
REQUIRED_VARS := ECR_ENDPOINT ECR_NAME ECS_SERVICE_NAME ECS_CLUSTER_NAME AWS_ACCESS_KEY AWS_SECRET_KEY AWS_REGION
$(foreach var,$(REQUIRED_VARS),$(if $(value $(var)),,$(error $(var) is undefined)))

# Get git commit hash for tagging
COMMIT_HASH := $(shell git rev-parse --short HEAD)
IMAGE_NAME := $(ECR_ENDPOINT)/$(ECR_NAME)

# AWS CLI configuration
export AWS_ACCESS_KEY_ID=$(AWS_ACCESS_KEY)
export AWS_SECRET_ACCESS_KEY=$(AWS_SECRET_KEY)
export AWS_DEFAULT_REGION=$(AWS_REGION)

.PHONY: all
all: build push deploy

.PHONY: build
build:
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME):$(COMMIT_HASH) .
	docker tag $(IMAGE_NAME):$(COMMIT_HASH) $(IMAGE_NAME):latest

.PHONY: push
push:
	@echo "Logging into ECR..."
	aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin $(ECR_ENDPOINT)
	
	@echo "Checking for existing latest tag..."
	if aws ecr-public describe-images --repository-name $(ECR_NAME) --image-ids imageTag=latest --region us-east-1 >/dev/null 2>&1; then \
		echo "Removing old latest tag..."; \
		aws ecr-public batch-delete-image --repository-name $(ECR_NAME) --image-ids "[{\"imageTag\":\"latest\"}]" --region us-east-1; \
	fi
	
	@echo "Pushing images to ECR..."
	docker push $(IMAGE_NAME):$(COMMIT_HASH)
	docker push $(IMAGE_NAME):latest

.PHONY: deploy
deploy:
	@echo "Updating ECS service..."
	aws ecs update-service \
		--cluster $(ECS_CLUSTER_NAME) \
		--service $(ECS_SERVICE_NAME) \
		--force-new-deployment

.PHONY: clean
clean:
	@echo "Cleaning up local images..."
	-docker rmi $(IMAGE_NAME):$(COMMIT_HASH)
	-docker rmi $(IMAGE_NAME):latest

.PHONY: check-env
check-env:
	@echo "Checking environment variables..."
	@echo "ECR Endpoint: $(ECR_ENDPOINT)"
	@echo "ECR Name: $(ECR_NAME)"
	@echo "ECS Service: $(ECS_SERVICE_NAME)"
	@echo "ECS Cluster: $(ECS_CLUSTER_NAME)"
	@echo "AWS Region: $(AWS_REGION)"
	@echo "Commit Hash: $(COMMIT_HASH)"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Build, push, and deploy"
	@echo "  build        - Build Docker image with commit hash and latest tags"
	@echo "  push         - Push images to ECR (removes old latest tag)"
	@echo "  deploy       - Deploy to ECS"
	@echo "  clean        - Remove local Docker images"
	@echo "  check-env    - Display current environment variables"
	@echo "  help         - Show this help message"