package main

import (
	"fmt"
	"os"
	"testing"
)

type TestReaderWriter struct {
	inner []byte
}

func (tw *TestReaderWriter) Write(p []byte) (n int, err error) {
	tw.inner = append(tw.inner, p...)
	return len(p), nil
}

func (tw *TestReaderWriter) Read(p []byte) (n int, err error) {
	n = copy(p, tw.inner)
	return n, nil
}

func Test_show_noNotes(t *testing.T) {
	notes := make(map[string]string)
	output := TestReaderWriter{}
	expected := "No notes yet\n"
	show(&output, notes)
	if result := string(output.inner); result != expected {
		t.Errorf("error")
	}
}
func Test_showNotes(t *testing.T) {
	notes := map[string]string{
		"notes1": "kjhdfkjd fdjkhgr",
	}
	output := TestReaderWriter{}
	show(&output, notes)
	expected := "notes1: 16\n"
	if result := string(output.inner); result != expected {
		t.Errorf("error")
	}
}
func Test_addNew(t *testing.T) {
	notes := make(map[string]string)
	output := TestReaderWriter{}
	noteTitle := "note1"
	input := TestReaderWriter{
		inner: []byte(fmt.Sprintf("%s\n", noteTitle)),
	}
	if err := addNew(&output, &input, notes, "touch"); err != nil {
		t.Errorf("add threw an error: %v", err)
	}
	if len(notes) != 1 {
		t.Error("add didn't add any notes")
	}
}

func Test_addNew_exit_editor(t *testing.T) {
	notes := make(map[string]string)
	output := TestReaderWriter{}
	noteTitle := "note1"
	input := TestReaderWriter{
		inner: []byte(fmt.Sprintf("%s\n", noteTitle)),
	}
	addNew(&output, &input, notes, "cat")
	if len(notes) != 1 {
		t.Error("add didn't add any notes")
	}
}

func Test_deleteFileName(t *testing.T) {
	noteTitle := "note1"
	notes := map[string]string{
		noteTitle: "uyyh hkfk kjy",
	}
	output := TestReaderWriter{}
	input := TestReaderWriter{
		inner: []byte(fmt.Sprintf("%s\n", noteTitle)),
	}
	deleteFileName(&output, &input, notes)

	if len(notes) != 0 {
		t.Error("add didn't add any notes")
	}
}

func Test_editFile(t *testing.T) {
	noteTitle := "note1"
	notes := map[string]string{
		noteTitle: " gsifgeiru ",
	}
	output := TestReaderWriter{}

	input := TestReaderWriter{
		inner: []byte(fmt.Sprintf("%s\n", noteTitle)),
	}
	if err := edit(&output, &input, notes, "cat"); err != nil {
		t.Errorf("edit threw an error, %v", err)
	}
}

func Test_edit_ChangedLength(t *testing.T) {
	noteTitle := "note1"
	notes := map[string]string{
		noteTitle: "note1\n",
	}
	output := TestReaderWriter{}
	input := TestReaderWriter{
		inner: []byte(fmt.Sprintf("%s\n", noteTitle)),
	}
	edit(&output, &input, notes, "sed", "-i", "", "-e", "s/note1/note123/g")
	newValue := notes[noteTitle]
	if len(newValue) == 5 {
		t.Errorf("edit didnt work")
	}
	expectedValue := "note123\n"
	if newValue != expectedValue {
		t.Errorf("new value wrong, expected: %q, got: %q", expectedValue, newValue)
	}
}

// check what happens when the env var is not set
func Test_getEditor_unset(t *testing.T) {
	os.Unsetenv(EDITOR_ENV_VAR)
	editor := getEditor()
	if editor != DEFAULT_EDITOR {
		t.Error("if empty, editor should be default")
	}
}

func Test_getEditor_empty_string(t *testing.T) {
	os.Setenv(EDITOR_ENV_VAR, " ")
	editor := getEditor()
	if editor != DEFAULT_EDITOR {
		t.Error("if empty, editor should be default")
	}
}

func Test_getEditor_kkfj(t *testing.T) {
	expectedValue := "iosjdi"
	os.Setenv(EDITOR_ENV_VAR, expectedValue)
	editor := getEditor()
	if editor != expectedValue {
		t.Errorf("expected: %q, got %q", expectedValue, editor)
	}
}
