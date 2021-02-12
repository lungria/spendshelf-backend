FROM --platform=${BUILDPLATFORM} golang:1.15.8-alpine3.13 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src
COPY . ./

RUN echo "Building for ${TARGETOS}/${TARGETARCH}"
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /spendshelf-backend

FROM alpine:3.13

WORKDIR /root/
COPY --from=builder /spendshelf-backend .

CMD ["./spendshelf-backend"]
