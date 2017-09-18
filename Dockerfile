FROM golang:latest

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/cmd/power
RUN go get -v -d -u . && go install -v .

CMD ["/go/bin/power"]
