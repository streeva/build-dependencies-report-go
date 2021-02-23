FROM golang:1.15-alpine3.12 as build-stage
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build

FROM alpine AS alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /build/build-dependencies-report /usr/bin/build-dependencies-report

ENTRYPOINT ["build-dependencies-report"]