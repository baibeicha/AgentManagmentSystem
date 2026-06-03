FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_NAME
RUN if [ -z "$SERVICE_NAME" ]; then echo "Build argument SERVICE_NAME is missing" && exit 1; fi

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/service ./cmd/${SERVICE_NAME}

FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/bin/service /app/service

RUN chmod +x /app/service

CMD ["/app/service"]