FROM golang:1.26-alpine

# Install essential build tools
RUN apk update && apk add --no-cache git make build-base openssh-client

# Set the working directory inside the container
WORKDIR /app

EXPOSE 4000
COPY . /app/

RUN go mod tidy

CMD [ "make", "run-dockerized" ]