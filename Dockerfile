FROM golang:1.17

MAINTAINER erdoganakdeniz
WORKDIR /go/app
COPY ./main.go
RUN go get -d -v
RUN go build -v
RUN echo $PATH
RUN ls
RUN pwd

CMD ["./go-mongodb-docker"]