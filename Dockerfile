FROM golang:1.19-alpine3.16 AS builder

RUN apk add --no-cache build-base wget git curl

ADD . /app

WORKDIR /app

RUN go build -o unri_fusioner -ldflags="-w -s" cmd/*.go

FROM alpine:3.16

COPY --from=builder /app/unri_fusioner /app/unri_fusioner

CMD ["/app/unri_fusioner"]

ENTRYPOINT ["/app/unri_fusioner"]

EXPOSE 3000
