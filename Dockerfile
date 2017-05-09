FROM golang:1.8-alpine
MAINTAINER Regner Blok-Andersen <shadowdf@gmail.com>

ADD . /go/src/app
WORKDIR /go/src/app
RUN go get

EXPOSE 8080
CMD ["go", "run", "main.go"]