ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
# COPY jsons/ ./jsons/
RUN go build -v -o /run-app .


FROM debian:bookworm

COPY --from=builder /run-app /usr/local/bin/
COPY --from=builder /usr/src/app/jsons /usr/local/jsons
COPY --from=builder /usr/src/app/config.properties /usr/local/config.properties

WORKDIR /usr/local
CMD ["run-app"]