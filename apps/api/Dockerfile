FROM golang:1.22-bookworm AS build
WORKDIR /app
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download
COPY apps/api .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o membrow-api ./cmd/api

FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*
COPY --from=build /app/membrow-api ./membrow-api
COPY apps/api/.env.example ./.env.example
ENV API_PORT=8080
EXPOSE 8080
CMD ["./membrow-api"]
