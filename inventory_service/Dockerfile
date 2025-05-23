FROM golang:1.23.0-alpine

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o inventory_service_app ./inventory_service/cmd

EXPOSE 50051
CMD ["./inventory_service_app"]
