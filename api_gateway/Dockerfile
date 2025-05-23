FROM golang:1.23.0-alpine

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api_gateway_app ./api_gateway

EXPOSE 8081
CMD ["./api_gateway_app"]
