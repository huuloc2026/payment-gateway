# syntax=docker/dockerfile:1
FROM golang:1.22 as builder
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o /processor ./cmd/processor

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /processor /processor
COPY .env.example /.env
EXPOSE 2113
ENTRYPOINT ["/processor"]
