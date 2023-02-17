#
# Observer, Signer, Statechain
#

#
# Build
#
FROM golang:1.13 AS build

WORKDIR /app

COPY . .

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go mod verify
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o $GOBIN/mock-binance ./cmd/mock-binance

#
# Main
#
FROM alpine

RUN apk add --update jq curl nginx && \
    rm -rf /var/cache/apk/*

# Copy the compiled binaires over.
COPY --from=build /go/bin/mock-binance /usr/bin/

EXPOSE 26660
CMD "/usr/bin/mock-binance"
