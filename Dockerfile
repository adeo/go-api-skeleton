# build
FROM golang:1-alpine as builder

RUN rm -rf /var/cache/apk/* && rm -rf /tmp/*
RUN apk update
RUN apk --no-cache add -U make git musl-dev gcc openssh-client

ARG GITHUB_SSH_PRIVATEKEY
ENV GITHUB_SSH_PRIVATEKEY=$GITHUB_SSH_PRIVATEKEY

WORKDIR /go/src/github.com/adeo/turbine-go-api-skeleton
COPY . /go/src/github.com/adeo/turbine-go-api-skeleton
RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no" > /root/.ssh/config && \
    git config --global url."ssh://git@github.com/".insteadOf "https://github.com/" && \
    eval $(ssh-agent -s) && echo "${GITHUB_SSH_PRIVATEKEY}" | ssh-add - > /dev/null && \
    make deps generate test build
RUN if [ "$(git status --porcelain)" ]; then echo "Git tree has changes. Make sure all the generated files was updated by doing 'make generate', and the commit all the changed files." && git status && exit 1; fi

# run
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl
COPY --from=builder /go/src/github.com/adeo/turbine-go-api-skeleton/turbine-go-api-skeleton .

RUN adduser -D -u 1001 runner
USER 1001

CMD ["/turbine-go-api-skeleton"]

