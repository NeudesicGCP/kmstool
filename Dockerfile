# Use an Alpine golang container for building
FROM golang:1.11-alpine3.8 AS builder

# Keep this build well away from GOPATH
RUN mkdir -p /tmp/kmstool
COPY . /tmp/kmstool/
WORKDIR /tmp/kmstool

# Building with CGO enabled; add dependencies that are not present in the base
# image
RUN apk add --update --no-cache git gcc musl-dev openssl curl ca-certificates && \
    go build

# The kmstool container
FROM alpine:3.8
LABEL maintainer="neugcp@neudesic.com"

COPY --from=builder /tmp/kmstool/kmstool /usr/local/bin/

# Add well-known CA certificates to make TLS support seamless
RUN apk add --no-cache ca-certificates

# Go networking needs a valid nsswitch configuration.
RUN [ ! -e /etc/nsswitch.conf ] && \
    echo 'hosts: files dns' > /etc/nsswitch.conf

# By default, execute `kmstool --help` as nobody.
# Override to encrypt/decrypt, and to write files with the correct ownership
USER nobody
ENTRYPOINT [ "/usr/local/bin/kmstool" ]
CMD [ "--help" ]
