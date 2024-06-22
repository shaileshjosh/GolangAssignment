package main

//packages are imported using import
import (
	"encoding/json"
	"fmt"
	"io" // io/ioutil moved to io
	"net/http"
	"strings"
)

const GameOverString string = "Game over!"
const YouWinString string = "Congratulations,you won!"
const WordConstant string = "racecar"

func getWord() string {
	res, err := http.Get("https://random-word-api.herokuapp.com/word?number=5")
	if err != nil {
		return WordConstant
	}
	body, err := io.ReadAll(res.Body) // get body from res
	res.Body.Close()
	if res.StatusCode > 299 { // if status code not 200
		return WordConstant
	}
	if err != nil {
		return WordConstant
	}
	var words []string // array of 5 words

	err = json.Unmarshal(body, &words)

	if err != nil {
		return WordConstant
	}
	return words[0] // return first word
}

func main() {

	word := getWord()

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
