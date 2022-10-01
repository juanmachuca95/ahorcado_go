# 💀 (Ahorcado) Multiplayer - Golang

Este juego surgue de la necesidad de experimentar en el maravilloso mundo de los servicios gRPC y sus correspondientes implementaciones. Soy novato en la creación de juegos multijugador, por favor disfruta del juego teniendo en cuenta sus limitaciones. 

Inspirado en el fantastico blog de 
 https://jbrandhorst.com/post/gopherjs-client-grpc-server-3/ Lectura sin desperdicio 👍

Probar el service con Evans con certificado tls: 
```zsh
evans --tls --cert cert/server-cert.pem --certkey cert/server-key.pem --cacert cert/ca-cert.pem --host "localhost" -r -p 8080 --verbose
```

Probar el service con grpcurl eludiendo el certificado:

```zsh
./grpcurl --insecure localhost:8080 list
```

Probar el service con curl:
```zsh
curl -k -X GET https://localhost:8080/api/v1/game 
```


Esta implementación permite jugarlo de manera local. 
Proximamente será subida una version online sencilla. 

Hecho con mucho cariño <b>@juanmachuca95</b>

## Seguir estudiando google cloud para implementar la api gRPC Gateway
https://cloud.google.com/api-gateway/docs/grpc-overview
https://medium.com/swlh/rest-over-grpc-with-grpc-gateway-for-go-9584bfcbb835
https://cloud.google.com/endpoints/docs/grpc/transcoding