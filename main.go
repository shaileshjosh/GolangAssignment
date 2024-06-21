package main

//packages are imported using import
import (
	"fmt"
	"strings"
)

const GameOverString string = "Game over!"
const YouWinString string = "Congratulations,you won!"

func main() {

	word := "racecar"

	// lookup for entries made by the user.
	entries := map[string]bool{}

	// list of "_" corrosponding to the number of letters in the word. [ _ _ _ _ _ ]
	placeholder := []string{}

	//create placeholder slice matching to length of word
	for i := 0; i < len(word); i++ {
		placeholder = append(placeholder, "_")
	}

	chances := 8

	for {

		// evaluate a loss! If user guesses a wrong letter or the wrong word, they lose a chance.
		userInput := strings.Join(placeholder, "")

		if chances == 0 && userInput != word {
			fmt.Println(GameOverString)
			break
		}
		// evaluate a win!
		if userInput == word {
			fmt.Println(YouWinString)
			break
		}

		//Console display
		fmt.Println()
		fmt.Println(placeholder)                 // render the placeholder
		fmt.Printf("Chances left:%d\n", chances) // render the chances left

		keys := []string{}
		for key, _ := range entries {
			keys = append(keys, key)
		}

		fmt.Println("Guesses: ", keys) //show the words/letters guessed till now.
		fmt.Printf("Guess the word or letter:")

		// take the input
		str := ""
		fmt.Scanln(&str)

		if len(str) > 1 { //check input is  word or single character
			if str == word {
				fmt.Println(YouWinString)
				break
			} else {
				entries[str] = true
				chances -= 1
				continue
			}
		}

		//compare and update entries ,placeholder and chances

		_, exist := entries[str]

		if exist {
			// already exist entry in the guesses
			continue
		}

		entries[str] = true
		found := false // flag to check character found

		//Iterate over the string to check character exists
		for i, v := range word {
			if str == string(v) {
				found = true
				placeholder[i] = string(v)
			}
		}
		//if not found,decrease the chances
		if !found {
			chances -= 1
		}

	}
}
