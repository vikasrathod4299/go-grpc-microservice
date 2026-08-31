FROM golang:1.26.1-apline AS builder

WORKDIR /src

copy go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOSE=linux go build \
  -trimpath \
  -ldflags="-s -w" \
  -o /out/gateway \
  ./cmd/gateway

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/gateway /gateway

EXPOSE 8080

ENTRYPOINT ["/gateway"]
