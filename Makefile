.PHONY: test

GOBIN?=${GOPATH}/bin

all: install

install: go.sum
	go install -v ./cmd/mock-binance

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	GO111MODULE=on go mod verify

run:
	mock-binance

test-coverage:
	@go test -v -coverprofile .testCoverage.txt ./...

coverage-report: test-coverage
	@tool cover -html=.testCoverage.txt

clear:
	clear

test:
	@go test -mod=readonly ./...

test-watch: clear
	@./scripts/watch.bash

lint: 
	@golangci-lint run
	@go mod verify

#=============== GITLAB =================#
docker-gitlab-login:
	docker login -u ${CI_REGISTRY_USER} -p ${CI_REGISTRY_PASSWORD} ${CI_REGISTRY}

docker-gitlab-push:
	docker push registry.gitlab.com/thorchain/bepswap/mock-binance:latest

docker-gitlab-build:
	docker build -t registry.gitlab.com/thorchain/bepswap/mock-binance .
	docker tag registry.gitlab.com/thorchain/bepswap/mock-binance $$(git rev-parse --short HEAD)
#========================================#
