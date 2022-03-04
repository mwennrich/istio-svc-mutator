FROM golang:1.17-alpine AS build
ENV GO111MODULE=on
ENV CGO_ENABLED=0


COPY . /app
WORKDIR /app

RUN go build istio-svc-mutator.go

FROM alpine:latest

WORKDIR /

COPY --from=build /app .

ENTRYPOINT ./istio-svc-mutator
