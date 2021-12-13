FROM golang:1.17

WORKDIR /src

RUN apt-get update && apt-get install -y \
    curl \
    netcat \
    && apt-get clean

RUN curl -o /usr/local/bin/whenavail https://bitbucket.org/silintl/docker-whenavail/raw/1.0.2/whenavail \
     && chmod a+x /usr/local/bin/whenavail

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

ADD ./application .

RUN go get ./...

RUN go build -o /bin/app server.go

RUN go build -o /bin/migrate migrations/main.go
RUN pwd && ls
CMD ["app", "dev"]
