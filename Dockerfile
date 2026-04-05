FROM golang:1.25-alpine AS build
WORKDIR /src

# Install sqlc
RUN apk add --no-cache curl && \
    curl -L https://github.com/sqlc-dev/sqlc/releases/download/v1.27.0/sqlc_1.27.0_linux_amd64.tar.gz | tar xz -C /tmp && \
    mv /tmp/sqlc /usr/local/bin/sqlc && \
    chmod +x /usr/local/bin/sqlc

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate sqlc code
RUN sqlc generate

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/api /api
EXPOSE 8080
USER nobody:nobody
ENTRYPOINT ["/api"]
