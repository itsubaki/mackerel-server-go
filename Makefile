SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)

runserver:
	-rm ${GOPATH}/bin/mackerel-api
	go install github.com/itsubaki/mackerel-api/cmd/mackerel-api
	GIN_MODE=debug mackerel-api
