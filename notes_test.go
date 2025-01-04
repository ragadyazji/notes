package main

import (
	//"bytes"
	//"io"
	"fmt"
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

func Test_addNew(t *testing.T) {
	notes := make(map[string]string)
	output := TestReaderWriter{}
	noteTitle := "note1"
	input := TestReaderWriter{
		inner: []byte(fmt.Sprintf("%s\n", noteTitle)),
	}
	addNew(&output, &input, notes, "touch")
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
