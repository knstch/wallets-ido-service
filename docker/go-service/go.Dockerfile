FROM golang:1.24 AS base

FROM base AS builder

WORKDIR /build
COPY . ./
RUN go build -mod=mod ./cmd/wallets

FROM base AS final

ARG PORT

WORKDIR /app
COPY --from=builder /build/wallets /build/.env ./
COPY --from=builder /build/wallets ./

EXPOSE ${PUBLIC_HTTP_ADDR}
CMD ["/app/wallets"]