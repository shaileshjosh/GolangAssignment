package main

//packages are imported using import
import (
	"encoding/json"
	"flag"
	"fmt"
	"io" // io/ioutil moved to io
	"net/http"
	"strings"
	"time"
)

const GameOverString string = "Game over!"
const YouWinString string = "Congratulations,you won!"
const WordConstant string = "racecar"
const MAX_CHANCES int = 8

// 1."dev" is flag 2.false is default value 3."dev mode" is helper text
var development = flag.Bool("dev", false, "dev mode")
var addr = flag.String("addr", "localhost:8085", "http service address")

type Hangman interface {
	RenderGame([]string, int, map[string]bool)
	getInput() string
}

type HangmanTerm struct {
}

func play(h Hangman, word string) bool {

	Entries := map[string]bool{}
	Placeholder := []string{}
	chances := MAX_CHANCES
	fmt.Println(chances)

	// create placeholder slice matching to length of word
	for i := 0; i < len(word); i++ {
		Placeholder = append(Placeholder, "_")
	}
	timer := time.NewTimer(10 * time.Minute)
	result := make(chan bool)

	go func() {
		for {

			// evaluate a loss! If user guesses a wrong letter or the wrong word, they lose a chance.
			userInput := strings.Join(Placeholder, "")

			if chances == 0 && userInput != word {
				result <- false
				return
			}
			// evaluate a win!
			if userInput == word {
				result <- true
				return
			}

			//Console display
			h.RenderGame(Placeholder, chances, Entries)

			// take the input
			str := h.getInput()

			if len(str) > 1 { //check input is  word or single character
				if str == word {
					result <- true
					return
				} else {
					Entries[str] = true
					chances -= 1
					continue
				}
			}

			//compare and update entries ,placeholder and chances

			_, exist := Entries[str]

			if exist {
				// already exist entry in the guesses
				continue
			}

			Entries[str] = true
			found := false // flag to check character found

			//Iterate over the string to check character exists
			for i, v := range word {
				if str == string(v) {
					found = true
					Placeholder[i] = string(v)
				}
			}
			//if not found,decrease the chances
			if !found {
				chances -= 1
			}

		}
	}()

	for {
		select {
		case r := <-result:
			if r {
				return true
			} else {
				return false
			}

		case <-timer.C:
			return false

		}
	}

}

func getWord() string {

	if *development {
		return WordConstant
	}

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

	flag.Parse() // parse the flag before use

	go webGame()

	hangmanGameStruct := HangmanTerm{}

	result := play(&hangmanGameStruct, getWord())
	if result {
		fmt.Println("You win!")
	} else {
		fmt.Println("You lose")
	}

}

func (h *HangmanTerm) RenderGame(placeholder []string, chances int, entries map[string]bool) {
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
}
func (h *HangmanTerm) getInput() string {
	str := ""
	fmt.Scanln(&str)
	return str
}
