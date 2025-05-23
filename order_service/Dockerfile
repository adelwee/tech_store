FROM golang:1.23.0-alpine

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order_service_app ./order_service/cmd

EXPOSE 50052
CMD ["./order_service_app"]
