FROM golang:latest as builder

WORKDIR /src
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /spendshelf-backend

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /spendshelf-backend .

CMD ["./spendshelf-backend"]