# Makefile for Figma Parser Backend Service

start: ## Run the application with Docker Compose (always build)
	cd docker && docker-compose up --build -d

stop: ## Stop Docker containers
	cd docker && docker-compose down -v

restart: ## Restart Docker containers
	cd docker && docker-compose restart

start-logs: ## Run the application with Docker Compose and show logs (always build)
	cd docker && docker-compose up --build
