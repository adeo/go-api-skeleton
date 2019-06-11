# build
FROM golang:1-alpine as builder

RUN rm -rf /var/cache/apk/* && rm -rf /tmp/*
RUN apk update
RUN apk --no-cache add -U make git musl-dev gcc

WORKDIR /go/src/github.com/adeo/turbine-go-api-skeleton
COPY . /go/src/github.com/adeo/turbine-go-api-skeleton
RUN make deps generate test build

# run
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
COPY --from=builder /go/src/github.com/adeo/turbine-go-api-skeleton/turbine-go-api-skeleton .

RUN adduser -D -u 1001 runner
USER 1001

CMD ["/turbine-go-api-skeleton"]

