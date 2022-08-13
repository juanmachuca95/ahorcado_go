FROM golang:1.18

RUN mkdir /app
ADD . /app
WORKDIR /app

EXPOSE 8080
EXPOSE 8090

RUN go build -o apiserver cmd/server/server.go

CMD [ "./apiserver" ]
