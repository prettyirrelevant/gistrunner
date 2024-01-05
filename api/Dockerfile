# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:1.21.4-bullseye as build

WORKDIR /opt/app

COPY . .

RUN go mod download && \
    CGO_ENABLED=0 go build -o /opt/app/gistrunner ./main.go


# -----------------------------------------------------------------------------
#  Final Stage
# -----------------------------------------------------------------------------
FROM alpine:latest as final

WORKDIR /opt/app

RUN apk -U upgrade && \
    apk add --no-cache dumb-init ca-certificates

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

COPY --from=build /opt/app/gistrunner /opt/app/gistrunner

CMD ["./gistrunner"]
