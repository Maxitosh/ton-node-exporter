FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /ton-node-exporter ./cmd/ton-node-exporter

FROM alpine:3.20

WORKDIR /

COPY --from=build-stage /ton-node-exporter /ton-node-exporter

ENV HTTP_SERVER_PORT 9100
EXPOSE $HTTP_SERVER_PORT

ENTRYPOINT ["/ton-node-exporter"]