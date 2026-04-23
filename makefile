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
# Create DB docker container
run-mysql:
	docker run --name trug-mysql -e MYSQL_ROOT_PASSWORD=$(TRUGPASS) -d -p 3306:3306 -v $(pwd)/internal/database:/var/lib/mysql mysql
# Connect to container
connect-mysql:
	docker exec -it trug-mysql mysql -u root -p