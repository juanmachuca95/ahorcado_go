# 💀 (Ahorcado) Multiplayer - Golang

Server developed in go using gRPC for communications that allows you to play the hangman game in a multiplayer way.
For this process it is necessary to previously have an account in mongo.cloud with a database for the game loaded. You just need two collections users and game. And set the credentials ```MONGODB_NAME``` and ```MONGODB_PASSWORD```

```zsh
docker run --env MONGODB_NAME=xxx --env MONGODB_PASSWORD=xxx -p 8080:8080 juanmachuca95/ahorcado:v1 -d
```

This is the game server, but the client for it is also developed.

Hecho con mucho cariño <b>@juanmachuca95</b>