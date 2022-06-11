package helpers

import (
	"strings"
)

func AlreadyFound(character string, encontrados []string) bool {
	var result bool = false
	for _, encontrado := range encontrados {
		if strings.ToUpper(character) == strings.ToUpper(encontrado) {
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