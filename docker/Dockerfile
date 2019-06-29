# build stage
FROM golang:1.12.6-alpine3.10 AS build-env

RUN apk update && \
    apk add --no-cache git ca-certificates

ADD . /work/

WORKDIR /work

ENV GO111MODULE=on

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /work/app
 
# final stage
FROM alpine:3.10

COPY --from=build-env /work/app /work/app
WORKDIR /work

EXPOSE 8080

CMD ["./app"]