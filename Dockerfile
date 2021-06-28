FROM golang:1.15.8-alpine3.13 as builder

WORKDIR /src
COPY . ./

RUN CGO_ENABLED=0 go build -o /spendshelf-backend

FROM alpine:3.13

WORKDIR /root/
COPY --from=builder /spendshelf-backend .

CMD ["./spendshelf-backend"]
