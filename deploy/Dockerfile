FROM golang:1.10.3-alpine3.7 as builder
RUN set -eux; \
    apk add --no-cache --virtual .build-dependencies git curl
WORKDIR $GOPATH/src/github.com/kunhou/TMDB
ADD . ./
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh  
RUN dep ensure -vendor-only
RUN go build -o /bin/movie ./cmd/movie

FROM alpine:latest
RUN set -eux; \
	apk add --no-cache ca-certificates openssl
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
COPY --from=builder /bin/movie ./
ENTRYPOINT [ "./movie" ]