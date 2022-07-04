FROM golang:1.18-alpine

WORKDIR $GOPATH/src/wbroker

COPY . .

RUN go get -d ./...
RUN go build -o /wbroker

EXPOSE 24005

CMD ["/wbroker"]