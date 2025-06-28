package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	// Validate required arguments
	if from == "" || to == "" {
		fmt.Println("Usage: go run . -from <source> -to <destination> [-offset <bytes>] [-limit <bytes>]")
		flag.PrintDefaults()
		return
	}

	// Call your Copy function
	err := Copy(from, to, offset, limit)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Copy completed successfully!")
	// Place your code here.
}
