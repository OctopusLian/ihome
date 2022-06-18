
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/userOrder/userOrder.proto

.PHONY: build
build: proto

	go build -o userOrder-srv main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t userOrder-srv:latest
