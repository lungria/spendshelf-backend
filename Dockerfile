FROM --platform=${BUILDPLATFORM} golang:1.15.3-alpine3.12 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src
COPY . ./

RUN echo "Building for ${TARGETARCH}/${TARGETARCH}"
RUN CGO_ENABLED=0 GOARCH=${TARGETARCH} GOOS=${TARGETOS} go build -o /spendshelf-backend

FROM alpine:3.12.1

WORKDIR /root/
COPY --from=builder /spendshelf-backend .

CMD ["./spendshelf-backend"]