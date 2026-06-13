FROM golang:1.25-alpine3.22 as builder

WORKDIR /kure

COPY go.mod .

RUN go mod download && go mod verify

RUN apk add --update --no-cache git

COPY . .

RUN CGO_ENABLED=0 go install -ldflags="-s -w" .

# ---------------------------------------------

FROM alpine:3.22
micro
RUN apk add --update --no-cache micro

COPY --from=builder /go/bin/sphinx /usr/bin/

CMD ["/usr/bin/kure"]
