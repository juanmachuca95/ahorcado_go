package frames

import "fmt"

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
		fmt.Print(s)
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
		fmt.Print(s)
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
		fmt.Print(s)
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
		fmt.Print(s)
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
		fmt.Print(s)
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
		fmt.Print(s)
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
		fmt.Print(s)
	}
}
