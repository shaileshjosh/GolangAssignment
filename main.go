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

// 1."dev" is flag 2.false is default value 3."dev mode" is helper text
var development = flag.Bool("dev", false, "dev mode")

type HangmanGame struct {
	Entries     map[string]bool
	Placeholder []string
	Chances     int
	Word        string
}

func play(h *HangmanGame, result chan bool) {

	// create placeholder slice matching to length of word
	for i := 0; i < len(h.Word); i++ {
		h.Placeholder = append(h.Placeholder, "_")
	}

	chances := 8

	for {

		// evaluate a loss! If user guesses a wrong letter or the wrong word, they lose a chance.
		userInput := strings.Join(h.Placeholder, "")

		if chances == 0 && userInput != h.Word {
			result <- false
			fmt.Println(GameOverString)
			return
		}
		// evaluate a win!
		if userInput == h.Word {
			result <- true
			fmt.Println(YouWinString)
			return
		}

		//Console display
		fmt.Println()
		fmt.Println(h.Placeholder)               // render the placeholder
		fmt.Printf("Chances left:%d\n", chances) // render the chances left

		keys := []string{}
		for key, _ := range h.Entries {
			keys = append(keys, key)
		}

		fmt.Println("Guesses: ", keys) //show the words/letters guessed till now.
		fmt.Printf("Guess the word or letter:")

		// take the input
		str := ""
		fmt.Scanln(&str)

		if len(str) > 1 { //check input is  word or single character
			if str == h.Word {
				result <- true
				fmt.Println(YouWinString)
				return
			} else {
				h.Entries[str] = true
				h.Chances -= 1
				continue
			}
		}

		//compare and update entries ,placeholder and chances

		_, exist := h.Entries[str]

		if exist {
			// already exist entry in the guesses
			continue
		}

		h.Entries[str] = true
		found := false // flag to check character found

		//Iterate over the string to check character exists
		for i, v := range h.Word {
			if str == string(v) {
				found = true
				h.Placeholder[i] = string(v)
			}
		}
		//if not found,decrease the chances
		if !found {
			chances -= 1
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

	hangmanGameStruct := HangmanGame{
		Entries:     map[string]bool{},
		Placeholder: []string{},
		Chances:     8,
	}
	hangmanGameStruct.Word = getWord()

	result := make(chan bool)

	timer := time.NewTimer(10 * time.Second)

	go play(&hangmanGameStruct, result)

	for {
		select {
		case <-result:
			fmt.Println("Game Ended")
			goto END
		case <-timer.C:
			fmt.Println("You have timedOut !!!")
			goto END

		}
	}
END:
	fmt.Println("Play Again..")

}
