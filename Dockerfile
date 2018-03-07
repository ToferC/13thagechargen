FROM golang:latest

ARG app_env
ENV APP_ENV $app_env

COPY . /go/src/github.com/toferc/13thAge/app
WORKDIR /go/src/github.com/toferc/13thAge/app

RUN go get ./
RUN go build

CMD go get github.com/pilu/fresh && fresh

EXPOSE 8080
