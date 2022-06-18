
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/getImg/getImg.proto

.PHONY: build
build: proto

	go build -o getImg-srv main.go plugin.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t getImg-srv:latest
