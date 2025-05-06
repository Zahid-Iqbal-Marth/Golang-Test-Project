FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o server server.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

ARG BATCH_SIZE=10
ARG BATCH_INTERVAL_SEC=60
ARG POST_ENDPOINT="https://eoa2dg5mkzbrlgw.m.pipedream.net"

ENV BATCH_SIZE=${BATCH_SIZE} \
    BATCH_INTERVAL_SEC=${BATCH_INTERVAL_SEC} \
    POST_ENDPOINT=${POST_ENDPOINT}

CMD ["./server"]