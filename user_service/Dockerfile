FROM golang:1.23.0-alpine

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user_service_app ./user_service/cmd

EXPOSE 50053
CMD ["./user_service_app"]
