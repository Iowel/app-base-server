# Stage 1
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
COPY ../../.env .env  

RUN go mod download

COPY ../../../../ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o app-server /app/cmd/app/main.go






# Stage 2
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/app-server /app/app-server
COPY --from=builder /app/.env /app/.env  

EXPOSE 8081

CMD ["/app/app-server"]


