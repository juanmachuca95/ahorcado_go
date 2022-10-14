# syntax=docker/dockerfile:1
FROM golang:1.19

WORKDIR /app

ARG name
ARG pass
ENV MONGODB_NAME $name
ENV MONGODB_PASSWORD $pass

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o myserver cmd/server/server.go

# Ports
EXPOSE 8080

CMD [ "./myserver" ]