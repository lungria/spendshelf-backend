FROM golang:1.13-alpine AS build-env
ENV GO111MODULE=on
RUN apk --no-cache add gcc build-base
WORKDIR /go/src/github.com/lungria/spendshelf-backend
ADD . .
RUN go mod tidy
RUN go build -o spendshelf-backend


FROM alpine:3.10
WORKDIR /go/bin
COPY --from=build-env /go/src/github.com/lungria/spendshelf-backend /go/bin
ENTRYPOINT ./spendshelf-backend