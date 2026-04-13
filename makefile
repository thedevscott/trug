# env variables can be used via the CLI interface ie:
# export TRUG_ADDR=":8889"
# go run ./cmd/web -addr=$TRUG_ADDR

help:
	go run ./cmd/web/ -help

run:
	go run ./cmd/web/ -addr=":4000"

web-test:
	go test -v ./cmd/web/

test:
	go test -v ./...

test-profile:
	go test -cover ./...
