package main

import (
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"gitlab.midas.dev/back/river/cmd"
)

func Ask4confirm() bool {
	var s string

	fmt.Printf("Pay a salary (y/N): ")
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func main() {
	if Ask4confirm() {
		cmd.Execute()
	}
	fmt.Println("Goodbye!")
}
