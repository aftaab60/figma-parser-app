# Makefile for Figma Parser Application

# Main application commands
start: ## Start the complete application (backend + frontend)
	@echo "Starting Figma Parser Application..."
	@echo "Starting Backend Services (Database + API)..."
	cd backend && make start
	@echo "Waiting for backend to be ready..."
	@sleep 3
	@echo "Starting Frontend Service..."
	cd frontend && make start
	@echo "Application started successfully!"
	@echo "Frontend: http://localhost:3001"
	@echo "Backend API: http://localhost:3000"
	@echo "Database: localhost:5432"

stop: ## Stop the complete application
	@echo "Stopping Figma Parser Application..."
	@echo "Stopping Frontend..."
	cd frontend && make stop || true
	@echo "Stopping Backend..."
	cd backend && make stop || true
	@echo "Application stopped successfully!"

restart: ## Restart the complete application
	@echo "Restarting Figma Parser Application..."
	make stop
	@sleep 5
	make start


start-backend: ## Start backend services only
	cd backend && make start

stop-backend: ## Stop backend services only
	cd backend && make stop

start-frontend: ## Start frontend service only
	cd frontend && make start

stop-frontend: ## Stop frontend service only
	cd frontend && make stop
