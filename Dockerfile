FROM golang:latest

RUN mkdir -p /go/src/go-docker-boilerplate
WORKDIR /go/src/go-docker-boilerplate

ADD . /go/src/go-docker-boilerplate

RUN go get -v

RUN  go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o chatbot" -command="./chatbot"