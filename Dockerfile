# FROM golang:1.8.3 as builder
#
# COPY . /go/src/github.com/maliceio/malice-api
# WORKDIR /go/src/github.com/maliceio/malice-api/
#
# RUN go get -u github.com/golang/dep/cmd/dep
# RUN dep ensure
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
#   -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o app .
#
# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
#
# WORKDIR /root/
#
# COPY --from=builder /go/src/github.com/maliceio/malice-api/app .
#
# CMD ["./app"]

FROM alpine:latest

LABEL maintainer "https://github.com/blacktop"

RUN apk --no-cache add ca-certificates
RUN apk --no-cache add -t .build-deps go git mercurial build-base \
  && git clone https://github.com/maliceio/malice-api.git /go/src/github.com/maliceio/malice-api \
  && cd /go/src/github.com/maliceio/malice-api/ \
  && export GOPATH=/go \
  && export PATH=/go/bin:$PATH \
  && go get -u github.com/golang/dep/cmd/dep \
  && dep ensure \
  && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/api \
  && rm -rf /go \
  && apk del --purge .build-deps

CMD ["api"]
