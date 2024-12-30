package main

import (
	// "bufio"
	// "errors"
	"fmt"
	// "os"
	// "strconv"
	// "strings"
)

var notes []string

func show() {
	for i, note := range notes {
		fmt.Println(i+1, note)
	}
}

func main() {
	var i int
	fmt.Println("select an option by the number: \n1. show all notes \n2. add a note\n3. delete a note \n4. edit  a note")
	fmt.Scan(&i)
	switch i {
	case 1:
		fmt.Println("You want to see your notes")
		show()
	// case 2:
	// 	fmt.Println("You want to add a note")
	// 	addError := add(os.Stdout, os.Stdin, tasks)
	// 	if addError != nil {
	// 		fmt.Fprintln(os.Stderr, addError)
	// 	}
	// case 3:
	// 	fmt.Println("You want to delete a note")
	// 	mark(os.Stdout, os.Stdin, tasks)
	// case 4:
	// 	fmt.Println("You want edit a note")
	default:
		fmt.Println("Choose between 1-4")
	}
}
