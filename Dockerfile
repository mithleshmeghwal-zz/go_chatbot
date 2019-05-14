FROM golang:latest

RUN mkdir -p /go/src/chatbot
WORKDIR /go/src/chatbot

ADD . /go/src/chatbot

RUN go get -v

RUN  go get github.com/githubnemo/CompileDaemon
RUN go get	github.com/gorilla/mux

ENTRYPOINT CompileDaemon -log-prefix=false -build="go build -o chatbot" -command="./chatbot"