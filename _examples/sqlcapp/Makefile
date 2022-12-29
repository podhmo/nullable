# need: go install github.com/kyleconroy/sqlc/cmd/sqlc@latest

gen:
	sqlc generate
.PHONY: gen

build:
	go get github.com/lib/pq
	go build ./...
.PHONY: build