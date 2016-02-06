FROM gliderlabs/alpine:3.3

MAINTAINER blacktop, https://github.com/blacktop

COPY . /var/www/malice-api
RUN apk-install nodejs git \
  && set -x \
	&& echo "Running npm install..." \
	&& cd /var/www/malice-api \
  && npm install \
  && npm install -g webpack-dev-server

WORKDIR /var/www/malice-api

ENTRYPOINT /usr/bin/webpack-dev-server --host 0.0.0.0
