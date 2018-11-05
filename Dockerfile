FROM golang:1.11

RUN apt-get -y update && apt-get -y install tzdata

WORKDIR /go/src/app
COPY . .

RUN go get github.com/mrbuk/dodns
RUN go install github.com/mrbuk/dodns
