FROM golang:1.16.6-alpine3.14 as builder

WORKDIR /src
COPY . ./

RUN go build -o /spendshelf-backend

FROM alpine3.14

WORKDIR /root/
COPY --from=builder /spendshelf-backend .

CMD ["./spendshelf-backend"]
