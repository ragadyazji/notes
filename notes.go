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
	for key, value := range notes {
		fmt.Printf("%s: %v\n", key, len(value))
	}
}

func FileChanged(initialStat os.FileInfo, file *os.File, c chan<- bool) {
	for {
		stat, err := file.Stat()
		if err != nil {
			return
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			c <- true
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func ProgramClosed(c chan<- bool, command *exec.Cmd) {
	command.Wait()
	// Wait for program termination signal
	c <- true
}

const EDITOR_FILENAME = "./RAGADSFILE"

func addNew(w io.Writer, r io.Reader, notes NoteList, editorCommand string, editorArgs ...string) {
	//print request for input for title
	reader := bufio.NewReader(r)
	fmt.Fprintln(w, "Enter note title: ")
	// read input and store in variable for title
	input, _ := reader.ReadString('\n')
	noteTitle := strings.TrimSpace(input)

	//make a temp file
	file, err := os.Create(EDITOR_FILENAME)
	if err != nil {
		log.Fatal(err)
	}
	if err := file.Sync(); err != nil {
		log.Fatal(err)
	}

	//print request for input for text
	fmt.Fprintln(w, "Enter lines of text:")
	contents, err := editFile(file, editorCommand, editorArgs)
	if err != nil {
		log.Fatal(err)
	}
	notes[noteTitle] = contents
}

// edits a temp file with the given editor and returns the contents or an error
func editFile(file *os.File, editorCommand string, editorArgs []string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// spawn a editor for the file
	cmd := exec.CommandContext(ctx, editorCommand, append(editorArgs, file.Name())...)
	// watch the file for edits
	initialStat, err := file.Stat()
	if err != nil {
		return "", err
	}

	cmd.Start()

	c := make(chan bool)
	go FileChanged(initialStat, file, c)
	go ProgramClosed(c, cmd)

	<-c

	// when saved copy the contents to the storage at the `title`
	contents, err := os.ReadFile(file.Name())
	if err != nil {
		return "", err
	}

	// clear the contents of the temp file (or destroy it)
	err = os.Remove(file.Name())
	if err != nil {
		return "", err
	}
	// TODO: close the spawned program
	return string(contents), nil
}

func deleteFileName(w io.Writer, r io.Reader, notes NoteList) {
	show(os.Stdout, notes)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the file name of the note you want to delete: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	delete(notes, input)
	// _, err := strconv.Atoi(input)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	show(os.Stdout, notes)
}

func edit(w io.Writer, r io.Reader, notes NoteList, editorCommand string, editorArgs ...string) {
	show(os.Stdout, notes)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the file name of the note you want to edit: ")
	input, _ := reader.ReadString('\n')
	noteTitle := strings.TrimSpace(input)
	value, ok := notes[noteTitle]
	if ok {
		//make a temp file
		file, err := os.Create(EDITOR_FILENAME)
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString(value)
		if err != nil {
			log.Fatal(err)
		}
		contents, err := editFile(file, editorCommand, editorArgs)
		if err != nil {
			log.Fatal(err)
		}
		notes[noteTitle] = contents
	} else {
		fmt.Println("Note not found")
	}
}

const DEFAULT_EDITOR = "code"
const DEFAULT_ARG = "--wait"

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
			addNew(os.Stdin, os.Stdout, notes, DEFAULT_EDITOR, DEFAULT_ARG)
		// 	// if err := add(); err != nil {
		// 	// 	fmt.Println("Error:", err)
		// 	// }
		case 3:
			fmt.Println("You want to delete a note")
			deleteFileName(os.Stdin, os.Stdout, notes)
		case 4:
			fmt.Println("You want edit a note")
			edit(os.Stdin, os.Stdout, notes, DEFAULT_EDITOR, DEFAULT_ARG)
		default:
			fmt.Println("Choose between 1-4")
		}
	}

}
