package main

//packages are imported using import
import (
	"fmt"
	"strings"
)

/*
Global variable should be declared with long declaration like "var",
short delcaration using := not supported
*/
var variable1 int = 5

func main() {
	fmt.Println("Value for variable1 is :", variable1)

	//array example
	arrayData := [4]string{"one", "two", "three", "Four"}
	fmt.Println("data in the array is :", arrayData)

	//slice example
	sliceData := []string{"a", "b", "c", "d"}
	fmt.Println("before appendind Slice Data is :", sliceData)

	//append is possible with slice
	sliceData = append(sliceData, "e")
	fmt.Println("after appendind Slice Data is :", sliceData)

	// for loop with iteration example
	for i := 0; i < len(sliceData); i++ {
		fmt.Println("value at", i, "is : ", sliceData[i])
	}

	//map example
	students := map[string]bool{"Surendra": true, "Vijay": false, "Saurabh": true}

	//example of length
	fmt.Println("Total number of the students are  ", len(students))

	//example of range and switch
	for key, value := range students {
		switch value {
		case true:
			fmt.Println(key, ": student is passed")
		case false:
			fmt.Println(key, ": student is failed")
		}
	}

	//hangman example

	fmt.Println()

	word := "Magnum"

	entries := map[string]bool{}

	placeholder := []string{}

	//create placeholder slice matching to length of word
	for i := 0; i < len(word); i++ {
		placeholder = append(placeholder, "_")
	}

	chances := 8

	for {

		concatedString := strings.Join(placeholder, "")

		if chances == 0 && concatedString != word {
			fmt.Println("Game over!")
			break
		}

		fmt.Println()
		fmt.Println(placeholder)

		if concatedString != word {
			fmt.Println("You have changes left :", chances)
		} else {
			fmt.Println("Congratulations,You won!")
			break
		}

		keys := []string{}
		for key, _ := range entries {
			keys = append(keys, key)
		}

		fmt.Println("Guesses: ", keys)

		fmt.Println("Guess the word or letter:")
		str := ""
		fmt.Scanln(&str)

		if !strings.Contains(word, str) {
			chances--
			entries[str] = false
		} else {
			entries[str] = true
		}
	}
}
