package main

import (
	"fmt"
	"os"
)

var (
	commitID string
)

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "version" {
			fmt.Printf("gophr %s\n", commitID)
		}
	}

	fmt.Println("Hello gophr!")
}