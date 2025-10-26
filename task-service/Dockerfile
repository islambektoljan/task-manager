FROM golang:1.25-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o /task-service

EXPOSE 8082

CMD ["/task-service"]