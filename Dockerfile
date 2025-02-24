FROM debian:stable-slim

# COPY source destination
COPY go.mod go.sum ./


COPY . .
CMD ["./web"]
