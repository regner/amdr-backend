FROM google/golang
MAINTAINER Regner Blok-Andersen <shadowdf@gmail.com>

ADD . /gopath/src/github.com/regner/amdr-backend
WORKDIR /gopath/src/github.com/regner/amdr-backend
RUN go get

EXPOSE 8080
CMD ["go", "run", "main.go"]