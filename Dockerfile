FROM golang:1.17-alpine AS builder
WORKDIR /go/src/
COPY . ./
RUN go build main.go

FROM alpine:3.14
COPY --from=builder /go/src/main /

ENTRYPOINT ["/main"]

