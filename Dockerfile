FROM golang:1.21.1

LABEL org.opencontainers.image.source='https://github.com/jacobbrewer1/f1-data'
LABEL org.opencontainers.image.description="A simple API that returns F1 data from the F1 Archive."
LABEL org.opencontainers.image.licenses='GNU General Public License v3.0'

WORKDIR /cmd

# Copy the binary from the build
COPY ./bin/app /cmd/app

RUN ["chmod", "+x", "./app"]

ENTRYPOINT ["/cmd/app"]
