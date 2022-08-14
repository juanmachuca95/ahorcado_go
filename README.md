# 💀 (Ahorcado) Multiplayer - Golang

Este juego surgue de la necesidad de experimentar en el maravilloso mundo de los servicios gRPC y sus correspondientes implementaciones. Soy novato en la creación de juegos multijugador, por favor disfruta del juego teniendo en cuenta sus limitaciones. 

Inspirado en el fantastico blog de 
 https://jbrandhorst.com/post/gopherjs-client-grpc-server-3/ Lectura sin desperdicio 👍

Esta implementación permite jugarlo de manera local. 
Proximamente será subida una version online sencilla. 


Hecho con mucho cariño <b>@juanmachuca95</b>

Implementación de certificado tls - Server GRPC - REST. Nota: la flag -insecure es necesaria

```bash
 ./grpcurl -insecure localhost:8080 protos.Ahorcado.GetGame
{
  "id": "62d41462b12b61ea29b74b5f",
  "word": "CORRUPTI",
  "encontrados": [
    "T",
    "C",
    "O"
  ]
}

```


