FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git build-base
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata dumb-init
WORKDIR /app

COPY --from=builder /app/app .

RUN adduser -D -g '' pruser
USER pruser

EXPOSE 8080

ENTRYPOINT ["dumb-init", "--"]
CMD ["./app"]
