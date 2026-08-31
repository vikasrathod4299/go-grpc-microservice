# syntax=docker/dockerfile:1

FROM golang:1.26.1-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/location \
    ./cmd/location

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/location /location

EXPOSE 50051

ENTRYPOINT ["/location"]
