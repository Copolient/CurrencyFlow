FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o currencyflow ./cmd/server/ && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o currencyflow-seed ./cmd/seed/

FROM alpine:3.20
RUN apk --no-cache add ca-certificates && \
    addgroup -S app -g 10001 && \
    adduser -S -G app -u 10001 app
WORKDIR /app
COPY --from=builder --chown=app:app /app/currencyflow .
COPY --from=builder --chown=app:app /app/currencyflow-seed .
COPY --from=builder --chown=app:app /app/config ./config

USER app
EXPOSE 3000
CMD ["./currencyflow"]
