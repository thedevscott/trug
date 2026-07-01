#===========================================================================
IMAGE_NAME = trug_go_local_dev
CONTAINER_NAME = trug_go_dev_env
WORK-DIR = /app

# Build the local Docker image using the Dockerfile
#===========================================================================
.PHONY: build
build:
	docker build -t $(IMAGE_NAME) .

# Run the container with the local directory mounted
#===========================================================================
.PHONY: dev
dev:
	@echo "Launching development environment..."
	docker run -it --rm \
		--name $(CONTAINER_NAME) \
		-p 4000:4000 \
		-v "$(PWD)":$(WORK-DIR) \
		-w $(WORK-DIR) \
		$(IMAGE_NAME) \
		sh

# Run Docker Compose up
#===========================================================================
.PHONY: compose-up
compose-up:
	@echo "Running docker compose up.."
	docker compose up --build

# Run Docker Compose down
#===========================================================================
.PHONY: compose-down
compose-down:
	@echo "Running docker compose down.."
	docker compose down -v

# Run Docker Compose down & Compose up
#===========================================================================
.PHONY: recompose
recompose: compose-down compose-up
	@echo "Cleaning up and re-running docker compose up.."

# Run Docker Compose db ogs
#===========================================================================
.PHONY: compose-db-logs
compose-db-logs:
	@echo "Running docker compose logs.."
	docker compose logs -f db

#===========================================================================
# env variables can be used via the CLI interface ie:
# export TRUG_ADDR=":8889"
# go run ./cmd/web -addr=$TRUG_ADDR

.PHONY: help
help:
	go run ./cmd/web/ -help
#===========================================================================

.PHONY: run
run:
	go run ./cmd/web/ -addr=":4000"
#===========================================================================

.PHONY: run-dockerized
run-dockerized:
	go run ./cmd/web/ -addr=":4000" \
	-dsn="web:trugpass@tcp(db:3306)/trug?parseTime=true"

#===========================================================================
.PHONY: web-test
web-test:
	go test -v ./cmd/web/

#===========================================================================
.PHONY: test
test:
	go test -v ./...

#===========================================================================
.PHONY: test-profile
test-profile:
	go test -cover ./...