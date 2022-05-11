FROM golang:1.18.1-alpine3.14 as builder

ARG TARGETOS
ARG TARGETARCH

ENV GOOS $TARGETOS
ENV GOARCH $TARGETARCH
ENV CGO_ENABLED 0

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ./

RUN go build -o /spendshelf-backend

FROM alpine:3.14

WORKDIR /root/
COPY --from=builder /spendshelf-backend .

CMD ["./spendshelf-backend"]
