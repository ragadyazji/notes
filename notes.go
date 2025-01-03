package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os/exec"
	"strings"
	"time"

	//"errors"
	"fmt"
	"os"
	//"strconv"
	//"strings"
)

type NoteList map[string]string

func show(w io.Writer, notes NoteList) {
	if len(notes) == 0 {
		fmt.Println("No notes yet")
		return
	}
	for key := range notes {
		fmt.Println(key)
	}
	//fmt.Println(notes)
}
func FileChanged(initialStat os.FileInfo, file *os.File) {
	for {
		stat, err := file.Stat()
		if err != nil {
			return
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

const EDITOR_FILENAME = "RAGADSFILE"

func addNew(w io.Writer, r io.Reader, notes NoteList, editorCommand string) {
	//print request for input for title
	reader := bufio.NewReader(r)
	fmt.Print("Enter note title: ")
	// read input and store in variable for title
	input, _ := reader.ReadString('\n')
	noteTitle := strings.TrimSpace(input)

	//make a temp file
	file, err := os.CreateTemp(os.TempDir(), EDITOR_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	// defer delete the file

	//print request for input for text
	fmt.Println("Enter lines of text:")
	ctx, _ := context.WithCancel(context.Background())
	// spawn a editor for the file
	cmd := exec.CommandContext(ctx, editorCommand, file.Name())
	// watch the file for edits
	// initialStat, err := file.Stat()
	// if err != nil {
	// 	return
	// }
	// go FileChanged(initialStat, file)

	cmd.Start()

	cmd.Wait()

	// when saved copy the contents to the storage at the `title`
	var contents []byte
	_, err = file.Read(contents)
	if err != nil {
		return
	}

	notes[noteTitle] = string(contents)

	//go ProgramClosed()

	// clear the contents of the temp file (or destroy it)
	// close the spawned program
}
func add(w io.Writer, r io.Reader, notes NoteList) {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter a title for the note: ")
	// title, _ := reader.ReadString('\n')
	// title = title[:len(title)-1]

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter file name: ")
	fileName, _ := reader.ReadString('\n')
	fileName = fileName[:len(fileName)-1] // Remove trailing newline

	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Get multiline input from user
	fmt.Println("Enter lines of text (empty line to finish):")
	var lines []string
	for {
		line, _ := reader.ReadString('\n')
		line = line[:len(line)-1] // Remove trailing newline

		if line == "" {
			break
		}
		lines = append(lines, line)
	}

	// Write lines to file
	for _, line := range lines {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// Create map with title as key and multiline text as value
	//textMap := make(map[string]string)
	notes[fileName] = ""
	for _, line := range lines {
		notes[fileName] += line + "\n"
	}
	show(w, notes)
}

// func deleteFileName() {
// 	show()
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter the file name of the note you want to delete: ")
// 	input, _ := reader.ReadString('\n')
// 	// _, err := strconv.Atoi(input)
// 	// if err != nil {
// 	// 	fmt.Println("Error:", err)
// 	// 	return
// 	// }

// 	delete(notes, input)
// 	show()
// }

// func edit() {

// }

func main() {
	notes := make(map[string]string)
	for {
		var i int
		fmt.Println("select an option by the number: \n1. show all notes \n2. add a note\n3. delete a note \n4. edit  a note")
		fmt.Scan(&i)
		switch i {
		case 1:
			fmt.Println("You want to see your notes")
			show(os.Stdout, notes)
		case 2:
			fmt.Println("You want to add a note")
			// add()
		// 	// if err := add(); err != nil {
		// 	// 	fmt.Println("Error:", err)
		// 	// }
		case 3:
			fmt.Println("You want to delete a note")
			// deleteFileName()
		// case 4:
		// 	fmt.Println("You want edit a note")
		default:
			fmt.Println("Choose between 1-4")
		}
	}

}
