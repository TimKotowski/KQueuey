FROM golang:1.23.2-bullseye as builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN go build -o kqueuey ./cmd

FROM debian:12

RUN useradd -d /home/kqueuey -u 777 kqueuey

COPY --from=builder /build/kqueuey /

USER kqueuey

CMD ["/kqueuey"]