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

RUN useradd user && mkdir /home/user && chown user.user /home/user && chown user.user /src
USER user
ENV GOPATH /home/user/go

# Copy the Go Modules manifests
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
COPY --chown=user application/go.mod go.mod
COPY --chown=user application/go.sum go.sum
RUN go mod download

COPY --chown=user ./application .

CMD ["air"]
