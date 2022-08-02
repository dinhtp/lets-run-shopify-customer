.PHONY: install test build serve clean pack deploy ship
TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)

export TAG

install:
	go get .

test: install
	go test ./...

build: install
	go build -ldflags "-X main.version=$(TAG)" -o ./bin/shopify-customer-service .

serve: build
	./bin/shopify-customer-service serve grpc

clean:
	rm -f ./bin/shopify-customer-service

dev:
	GOOS=linux make build
	docker build -t localhost:5000/shopify-customer-service .
	docker push localhost:5000/shopify-customer-service