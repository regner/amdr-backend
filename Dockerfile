FROM golang
MAINTAINER Regner Blok-Andersen <shadowdf@gmail.com>

ADD . /go/src/github.com/regner/amdr-backend
WORKDIR /go/src/github.com/regner/amdr-backend
RUN go get
RUN go install
ENTRYPOINT /go/bin/amdr-backend

EXPOSE 8080