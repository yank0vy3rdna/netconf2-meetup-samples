FROM ghcr.io/yank0vy3rdna/netconf2-meetup-samples:master

COPY --from=golang:1.22-alpine /usr/local/go/ /usr/local/go/

ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

ADD go.mod .

ADD go.sum .

ENV CGO_ENABLED 1

RUN go mod download

ADD . .

RUN go build -o app ./cmd

RUN go test ./...
