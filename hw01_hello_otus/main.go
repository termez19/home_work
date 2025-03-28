package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	reverseFunc := reverse.String("Hello, OTUS!")
	fmt.Println(reverseFunc)
}
