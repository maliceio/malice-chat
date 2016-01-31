FROM gliderlabs/alpine:3.3

MAINTAINER blacktop, https://github.com/blacktop

COPY . /go/src/github.com/maliceio/malice-api
# RUN apk-install
RUN apk-install -t build-deps go git mercurial \
  && set -x \
	&& echo "Building info Go binary..." \
	&& cd /go/src/github.com/maliceio/malice-api \
	&& export GOPATH=/go \
	&& go get \
	&& go build -ldflags "-X main.Version=$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/api \
	&& rm -rf /go \
	&& apk del --purge build-deps

ENTRYPOINT ["/bin/api"]

# CMD ["--help"]
