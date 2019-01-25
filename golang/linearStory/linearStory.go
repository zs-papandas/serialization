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

func playStory(page *storyPage){
	if page == nil{
		return
	}

	fmt.Println(page.text)
	page.playStory(page.nextPage)
}

func main()  {

	//scanner := bufio.NewScanner(os.Stdin)

	page1 := storyPage{"Line 01 for linked list", nil} 
	page2 := storyPage{"Line 02 for linked list", nil} 
	page3 := storyPage{"Line 03 for linked list", nil} 
	page1.nextPage = &page2
	page2.nextPage = &page3

	playStory(&page1)

	fmt.Println("Over")
}