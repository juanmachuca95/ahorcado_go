# (💀 Ahorcado) - Golang 

Esta es mi implementación del juego del ahorcado en Golang.

* Mongodb
* gRPC Server - Client 
* Stream bidirectional 
* Golang
* Game

Si te gusto el juego regalame un estrella. ⭐
Made by <b>@juanmachuca95</b>

### ¿Como juego?
Si quieres jugar de forma local, debes crear un archivo ```.env``` con tus variables con respecto a la base de datos en mongodb.

El gRPC Server se ejecuta:
```go
go run cmd/main/main.go
```

```go
go run client/client.go -addr=xxxxxxxxx:8080
```


