.PHONY: migrate
migrate: ## Run the migrations
	atlas migrate apply --env gorm

.PHONY: rollback
rollback: ## Rollback the migrations
	atlas migrate down --env gorm

.PHONY: generate-migration
generate-migration: ## Generate a new migration
	@printf "\033[33mEnter migration message: \033[0m"
	@read -r message; \
	atlas migrate diff --env gorm "$$message"