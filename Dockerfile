FROM golang
MAINTAINER Regner Blok-Andersen <shadowdf@gmail.com>

ADD . /go/src/github.com/regner/amdr-backend
RUN go install github.com/regner/amdr-backend
ENTRYPOINT /go/bin/amdr_backend

EXPOSE 8080