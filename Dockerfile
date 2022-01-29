FROM golang:1.17

WORKDIR /src

RUN apt-get update && apt-get install -y \
    curl \
    netcat \
    && apt-get clean

RUN curl -o /usr/local/bin/whenavail https://bitbucket.org/silintl/docker-whenavail/raw/1.0.2/whenavail \
     && chmod a+x /usr/local/bin/whenavail

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

ADD ./application .

RUN go get ./...

CMD ["air"]
