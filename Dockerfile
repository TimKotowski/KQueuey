FROM golang:1.23.2-bullseye as builder

WORKDIR /kqueuey

COPY go.mod go.sum

COPY . .

ENV CGO_ENABLED=0

RUN go build -o /kqueuey .

FROM gcr.io/distroless/static-debian12

COPY --from=builder /kqueuey .

CMD ["./kqueuey"]