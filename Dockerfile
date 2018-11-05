FROM golang:1.11

RUN apt-get -y update && apt-get -y install tzdata

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download github.com/mrbuk/dodns
RUN go-wrapper install github.com/mrbuk/dodns
