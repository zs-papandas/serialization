package main

import (
	"bufio"
	"fmt"
	"os"
)

// LINKED LIST

type storyPage struct {
	text string
	nextPage *storyPage
}

func (page *storyPage) playStory() {
	for page != nil {
		fmt.Println(page.text)
		page = page.nextPage
	}
}

func (page *storyPage) addToEnd(text string) {
	for page.nextPage != nil {
		page = page.nextPage
	}
	page.nextPage = &storyPage{text, nil}
}

func main()  {

	scanner := bufio.NewScanner(os.Stdin)

	page1 := storyPage{"Line 01 for linked list", nil} 
	page1.addToEnd("Line 02 for linked list")
	page1.addToEnd("Line 03 for linked list")
	
	
	for{
		fmt.Println("\n============================================")
		fmt.Println("Choose your options:")
		fmt.Println("(a) Add a new line")
		fmt.Println("(b) Play all the lines")
		fmt.Println("(c) Exit")

		scanner.Scan()
		response := scanner.Text()

		if response == "a" {
			fmt.Println("Type your text and press ENTER:")
			scanner.Scan()
			text := scanner.Text()
			page1.addToEnd(text)
		} else if response == "b" {
			page1.playStory()
		} else if response == "c" {
			fmt.Println("**** Application Over *****")
			break
		} else {
			fmt.Println("**** Invalid response, try again. ****")
		}
	}

}