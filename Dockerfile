FROM golang:1.23.2-bullseye as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

RUN go build -o ./bin/kqueuey ./cmd

FROM debian:12

RUN useradd -m -u 777 kqueuey

COPY --from=builder /build/bin/kqueuey /

USER kqueuey

ENTRYPOINT ["/kqueuey"]