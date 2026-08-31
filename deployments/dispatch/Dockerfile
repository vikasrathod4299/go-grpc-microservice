# syntax=docker/dockerfile:1

FROM golang:1.26.1-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/dispatch \
    ./cmd/dispatch

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/dispatch /dispatch

EXPOSE 50052

ENTRYPOINT ["/dispatch"]
