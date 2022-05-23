FROM golang:1.18

RUN mkdir /app
ADD . /app
WORKDIR /app

EXPOSE 8080

RUN go build -o apiserver server/server.go

CMD [ "./apiserver" ]
