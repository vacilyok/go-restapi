##
## Build
##
FROM golang:1.16-alpine AS build_base

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN apk add git

ARG GITLAB_USERNAME
ARG GITLAB_TOKEN

RUN echo "machine gitlab.net login ${GITLAB_USERNAME} password ${GITLAB_TOKEN}" > ${HOME}/.netrc

ENV EXPORT GO111MODULE=on
ENV GOPRIVATE=gitlab.net
RUN go mod download

COPY . ./
RUN go build -o /app/mediator .

##
## Deploy
##
FROM alpine:latest

WORKDIR /app

COPY --from=build_base /app .
COPY --from=build_base /build/migrations ./
COPY --from=build_base /build/scripts ./
EXPOSE 4000
CMD ["/app/mediator"]
