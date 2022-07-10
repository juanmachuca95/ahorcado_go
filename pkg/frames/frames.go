package frames

import (
	"github.com/pterm/pterm"
)

func Frames(try int) {
	switch try {
	case 0:
		Inicial(try)
	case 1:
		FirstAttempt(try)
	case 2:
		SecondAttempt(try)
	case 3:
		ThirdAttempt(try)
	case 4:
		FourthAttempt(try)
	case 5:
		FifthAttempt(try)
	case 6:
		SixthAttempt(try)
	default:
		Inicial(0)
	}
}

func Inicial(i int) {
	if i == 0 {
		s := `
		           _____
			   |   |
			  _O/  |
			   |   |
			  / \  |
			 ______|
		`
		pterm.FgRed.Println(s)
	}
}

func FirstAttempt(i int) {
	if i == 1 {
		s := `
		           _____
			   |   |
			  _O/  |
			   |   |
			  /    |
			 ______|
		`
		pterm.FgRed.Println(s)
	}
}

func SecondAttempt(i int) {
	if i == 2 {
		s := `
		           _____
			   |   |
			  _O/  |
			   |   |
			       |
			 ______|
		`
		pterm.FgRed.Println(s)
	}
}

func ThirdAttempt(i int) {
	if i == 3 {
		s := `
		           _____
			   |   |
			  _O/  |
			       |
			       |
			 ______|
		`
		pterm.FgRed.Println(s)
	}
}

func FourthAttempt(i int) {
	if i == 4 {
		s := `
		           _____
			   |   |
			  _O   |
			       |
			       |
			 ______|
		`
		pterm.FgRed.Println(s)
	}
}

func FifthAttempt(i int) {
	if i == 5 {
		s := `
		           _____
			   |   |
			   O   |
			       |
			       |
			 ______|
		`
		pterm.FgRed.Println(s)
	}
}

func SixthAttempt(i int) {
	if i == 6 {
		s := `
		           _____
			   |   |
			       |
			       |
			       |
			 ______|
		`
		pterm.FgRed.Println(s)
	}
}
