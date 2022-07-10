package helpers

import "testing"

func TestAlreadyFound(t *testing.T) {
	tt := struct {
		name        string
		encontrados []string
		word        string
		want        bool
	}{
		name:        "Testing Already Found character or word.",
		encontrados: []string{"A", "B", "C", "E"},
		word:        "e",
		want:        true,
	}

	if got := AlreadyFound(tt.word, tt.encontrados); got != tt.want {
		t.Log(tt.name)
		t.Errorf("Failed to test AlreadyFound, it's expected (%v), but was get (%v)", tt.want, got)
	}
}

func TestWin(t *testing.T) {
	tt := struct {
		name        string
		encontrados []string
		clave       string
		want        bool
	}{
		name:        "Testing Win last character is founded.",
		encontrados: []string{"M", "r", "C", "o", "s", "a"},
		clave:       "Marcos",
		want:        true,
	}

	if got := Win(tt.clave, tt.encontrados); got != tt.want {
		t.Log(tt.name)
		t.Errorf("Failed to test Win, it's expected (%v), but was get (%v)", tt.want, got)
	}
}

func TestShowWord(t *testing.T) {
	tt := struct {
		name   string
		letras []string
		clave  string
		want   string
	}{
		name:   "Testing Show Word",
		letras: []string{"P", "E", "M"},
		clave:  "PROBLEM",
		want:   "P____EM",
	}

	if got := ShowWord(tt.clave, tt.letras); got != tt.want {
		t.Log(tt.name)
		t.Errorf("Failed to test ShowWord, it's expected (%v), but was get (%v)", tt.want, got)
	}
}
