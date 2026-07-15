LOCAL_PATH := github.com/antongolenev23/voltake-services



# PRE-COMMIT HOOKS
.PHONY: fmt
fmt:
	@echo "Formatting all Go files..."
	@goimports -local $(LOCAL_PATH) -w . > /dev/null 2>&1

.PHONY: fmt-check
fmt-check:
	@echo "Checking all Go files format..."
	@if goimports -local $(LOCAL_PATH) -l . 2>&1 | grep . > /dev/null; then \
		echo "Not all Go files are formatted. Run 'make fmt'"; \
		exit 1; \
	else \
		echo "All Go files formatted correctly"; \
	fi

.PHONY: pre-commit
pre-commit: fmt-check
	@echo "Pre-commit checks passed"



# AUTH TESTS
COMPOSE=docker compose \
	-f deploy/compose/docker-compose.local.yml \
	-f deploy/compose/docker-compose.test.yml \
	--env-file .env

.PHONY: test-auth

test-auth:
	$(COMPOSE) down 

	$(COMPOSE) up auth-db auth-migrate auth-service -d --wait --build

	go test ./services/auth/tests/... -count=1 -parallel=16 -v

	$(COMPOSE) down 



# ALL SERVICES CONTROL LOCAL
COMPOSE_FILE=deploy/compose/docker-compose.local.yml
ENV_FILE=.env

.PHONY: up rebuild down logs

up:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up -d

rebuild:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) up --build -d

down:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) down

logs:
	docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) logs -f