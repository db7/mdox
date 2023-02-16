package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Writer struct {
	data []string
}

func (w *Writer) Print(args ...any) {
	w.data = append(w.data, fmt.Sprint(args...))
}
func (w *Writer) Printf(format string, args ...any) {
	w.data = append(w.data, fmt.Sprintf(format, args...))
}
func (w *Writer) Println(args ...any) {
	w.data = append(w.data, fmt.Sprintln(args...))
}

func (w *Writer) Fwrite(fd io.Writer) error {
	for _, d := range w.data {
		if _, err := fmt.Fprint(fd, d); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) Write(fn string) error {
	dir := filepath.Dir(fn)
	if err := os.MkdirAll(dir, 0775); err != nil {
		return err
	}

	fd, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer fd.Close()

	for _, d := range w.data {
		if _, err := fmt.Fprint(fd, d); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) Last() *string {
	if w.data == nil {
		return nil
	}
	return &w.data[len(w.data)-1]
}
