# Stage 1
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY ../../go.mod ../../go.sum ./ 
COPY ../../.env .env

COPY ./static /app/static

RUN go mod download

COPY ../../../../ ./ 

RUN CGO_ENABLED=0 GOOS=linux go build -o web-server /app/cmd/web/.  # Собираем бинарник




# Stage 2
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/web-server /app/web-server

COPY --from=builder /app/.env /app/.env

EXPOSE 8082

CMD ["/app/web-server"]
