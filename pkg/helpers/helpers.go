package helpers

import (
	"fmt"
	"strings"
)

const (
	_codeFound        = 1
	_codeNotFound     = 2
	_codeAlreadyFound = 3
	_codeWinner       = 4
	_codeUnexpected   = 5
)

func AlreadyFound(character string, encontrados []string) bool {
	var result bool = false
	for _, encontrado := range encontrados {
		if strings.EqualFold(character, encontrado) {
			result = true
		}
	}
	return result
}

func Win(clave string, encontrados []string) bool {
	var winner bool = false
	var claveLength int = len(clave)
	var encontradosLength int = 0

	for _, encontrado := range encontrados {
		count := strings.Count(strings.ToUpper(clave), strings.ToUpper(encontrado))
		encontradosLength += count
	}

	if encontradosLength == claveLength {
		winner = true
	}

	return winner
}

func ShowWord(clave string, letras []string) string {
	chars := []rune(clave)
	var wordPositions []string

	var encontrados []string = letras
	var results []string
	for i := 0; i < len(chars); i++ {
		char := string(chars[i])
		results = append(results, "_")
		wordPositions = append(wordPositions, char)
	}

	for _, encontrado := range encontrados {
		for i, worPos := range wordPositions {
			if encontrado == worPos {
				results[i] = worPos
			}
		}
	}

	return strings.Join(results, "")
}

func MessageStatus(user, myuser, word, winWord string, status int) (string, int) {
	messageMe := false
	if user == myuser {
		messageMe = true
	}
	switch status {
	case _codeFound:
		if messageMe {
			return fmt.Sprintf("¡Has encontrado la letra %s! 👍", word), 0
		}
		return fmt.Sprintf("El usario %s ha encontrado la letra %s 🔥", user, word), 0
	case _codeNotFound:
		if messageMe {
			return fmt.Sprintf("La letra %s no tiene coincidencias! 👎", word), 1
		}
		return fmt.Sprintf("¡El usuario %s ha intentado la letra %s sin exito! 🖊️", user, word), 0
	case _codeAlreadyFound:
		if messageMe {
			return fmt.Sprintf("La letra %s ya ha sido encontrada previamente! 🐢", word), 1
		}
		return fmt.Sprintf("¡El usuario %s ha intentado la letra %s pero ya ha sido encontrada! 🐢", user, word), 0
	case _codeWinner:
		if messageMe {
			return fmt.Sprintf("🎉 Has ganado! 🎊 La palabra era %s 🏆", winWord), 0
		}
		return fmt.Sprintf("¡El usuario %s ha ganado! 🐢 La palabra era: %s", user, winWord), 0
	case _codeUnexpected:
		return "Ha ocurrido un error inesperado", 0
	default:
		return "", 0
	}
}
