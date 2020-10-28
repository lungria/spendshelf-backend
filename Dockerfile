FROM golang:1.15.3-alpine3.12 as builder

WORKDIR /src
COPY . ./

RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o /spendshelf-backend

FROM alpine:3.12.1

WORKDIR /root/
COPY --from=builder /spendshelf-backend .

CMD ["./spendshelf-backend"]