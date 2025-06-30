.PHONY: migrate migrate_down migrate_up migrate_version docker prod docker_delve local swaggo test

LIST_GO_FILES = Get-ChildItem -Path . -Recurse -Filter *.go | ForEach-Object { $_.FullName }

# ==============================================================================
# Go migrate postgresql

force:
	migrate -database postgres://postgres:lalapopo123@localhost:5432/user_service_db?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://postgres:lalapopo123@localhost:5432/user_service_db?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://postgres:lalapopo123@localhost:5432/user_service_db?sslmode=disable -path migrations up

migrate_down:
	migrate -database postgres://postgres:lalapopo123@localhost:5432/user_service_db?sslmode=disable -path migrations down 1


# ==============================================================================
# Docker compose commands

develop:
	echo "Starting docker environment"
	docker-compose -f docker-compose.dev.yml up --build

docker_delve:
	echo "Starting docker debug environment"
	docker-compose -f docker-compose.delve.yml up --build

prod:
	echo "Starting docker prod environment"
	docker-compose -f docker-compose.prod.yml up --build

local:
	echo "Starting local environment"
	docker-compose -f docker-compose.yml up --build


# ==============================================================================
# Tools commands

run-linter:
	echo "Starting linters"
	golangci-lint run ./...

swaggo:
	echo "Starting swagger generating"
	swag init -g ./cmd/api/main.go

swaggo-windows:
	powershell -Command "{$oFiles = $(LIST_GO_FILES) -join ','; swag init -g $oFiles}"

# ==============================================================================
# Main

run:
	go run ./cmd/api/main.go

build:
	go build ./cmd/api/main.go

test:
	go test -cover ./...


# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache


# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)