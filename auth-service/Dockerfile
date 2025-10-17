FROM golang:1.25-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o /auth-service

EXPOSE 8081

CMD ["/auth-service"]