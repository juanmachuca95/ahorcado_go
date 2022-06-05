package game

import (
	database "github.com/juanmachuca95/ahorcado_go/internal/database/mongo"
)

var (
	db          = database.Connect()
	gameService = NewGameStorageGateway(db)
)

/* func TestInGame(t *testing.T) {
	tt := struct {
		name string
		word string
		user string
		id   string
		want string
	}{
		name: "Testing InGame character or word.",
		word: "a",
		user: "Juan",
		id:   "asdfas123123412d",
		want: "",
	}

	if got, err := gameService.inGame(tt.word, tt.user, tt.id); got != tt.want && err != nil {

	}
}
*/
