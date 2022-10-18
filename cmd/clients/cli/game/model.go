package game

import ah "github.com/juanmachuca95/ahorcado_go/protos/ahorcado"

type GameAhorcado struct {
	Id          string
	Encontrados []string
	Usersend    string
	Wordsend    string
	Winner      string
	Word        string
	Error       string
	Finalizada  bool
	Status      int64
}

func NewAhorcado(g *ah.Game) *GameAhorcado {
	return &GameAhorcado{
		Id:          g.Id,
		Encontrados: g.Encontrados,
		Usersend:    g.Usersend,
		Wordsend:    g.Wordsend,
		Word:        g.Word,
		Finalizada:  g.Finalizada,
		Status:      int64(g.Status),
		Error:       g.Error,
	}
}
