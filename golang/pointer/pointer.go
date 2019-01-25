package main

import "fmt"

func main(){
	x := 5
	fmt.Println(x)

	var xPtr *int = &x
	fmt.Println(xPtr)
}