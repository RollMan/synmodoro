FROM golang:1.14.7-alpine3.12 as develop
WORKDIR /work
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN apk add --no-cache git
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app -ldflags '-s -w'

from scratch as prod
WORKDIR /goapp
COPY --from=develop /work/app /goapp/app
ENTRYPOINT ['/goapp/app']
