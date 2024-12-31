package main

import (
	"bufio"
	//"errors"
	"fmt"
	"os"
	//"strconv"
	//"strings"
)

var notes map[string]string

func show() {
	if len(notes) == 0 {
		fmt.Println("No notes yet")
		return
	}
	// for key := range notes {
	// 	fmt.Println(key)
	// }
	fmt.Println(notes)
}

func add() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a title for the note: ")
	title, _ := reader.ReadString('\n')
	title = title[:len(title)-1]

	file, err := os.Create(title)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	fmt.Println("Enter text (Empty line to finish):")
	var lines []string
	for {
		line, _ := reader.ReadString('\n')
		line = line[:len(line)-1] // Remove trailing newline

		if line == "" {
			break
		}
		lines = append(lines, line)
	}

	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	notes[title] = ""
	for _, line := range lines {
		notes[title] += line + "\n"
	}
	show()
}

// func delete() {
// 	show()
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter the index of the note you want to delete: ")
// 	input, _ := reader.ReadString('\n')
// 	input = input[:len(input)-1]
// 	index, err := strconv.Atoi(input)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	if index < 1 || index > len(notes) {
// 		fmt.Println("Invalid index")
// 		return
// 	}
// 	notes = append(notes[:index-1], notes[index:]...)
// 	show()
// }

// func edit() {

// }

func main() {
	notes = make(map[string]string)
	for {
		var i int
		fmt.Println("select an option by the number: \n1. show all notes \n2. add a note\n3. delete a note \n4. edit  a note")
		fmt.Scan(&i)
		switch i {
		case 1:
			fmt.Println("You want to see your notes")
			show()
		case 2:
			fmt.Println("You want to add a note")
			add()
		// 	// if err := add(); err != nil {
		// 	// 	fmt.Println("Error:", err)
		// 	// }
		// case 3:
		// 	fmt.Println("You want to delete a note")
		// 	delete()
		// case 4:
		// 	fmt.Println("You want edit a note")
		default:
			fmt.Println("Choose between 1-4")
		}
	}

}
