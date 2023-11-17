FROM golang:1.20

WORKDIR /src

# silintl/docker-whenavail is little script that uses netcat to wait for a service to be ready
RUN apt-get update && apt-get install -y curl netcat-traditional && apt-get clean
RUN curl -o /usr/local/bin/whenavail https://bitbucket.org/silintl/docker-whenavail/raw/1.0.2/whenavail \
     && chmod a+x /usr/local/bin/whenavail

# cosmtrk/air is a project auto-build tool
RUN go install github.com/cosmtrek/air@v1.43.0

# pressly/goose is a database migrations tool
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.11.2

# set up to run as a normal user
RUN useradd user && mkdir /home/user && chown user:user /home/user && chown user:user /src
USER user
ENV GOPATH /home/user/go

# Copy the Go Modules manifests
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
COPY --chown=user ./go.mod go.mod
COPY --chown=user ./go.sum go.sum
RUN go mod download

COPY --chown=user ./ .

CMD ["air"]
