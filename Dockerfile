FROM golang:1.8-alpine
MAINTAINER Regner Blok-Andersen <shadowdf@gmail.com>

ADD . /go/src/github.com/regner/amdr-backend
WORKDIR /go/src/github.com/regner/amdr-backend
RUN go get

EXPOSE 8080
CMD ["go", "run", "main.go"]