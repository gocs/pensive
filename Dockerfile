FROM golang:alpine3.12 as build

WORKDIR /src
COPY . /src

RUN go build -o app ./cmd/main.go

EXPOSE 8080
ENTRYPOINT [ "./app" ]