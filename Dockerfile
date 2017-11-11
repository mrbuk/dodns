FROM golang:1.8

WORKDIR /go/src/app
COPY . .

RUN go-wrapper download github.com/mrbuk/dodns
RUN go-wrapper install github.com/mrbuk/dodns
