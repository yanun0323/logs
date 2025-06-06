package logs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Writer interface {
	io.Writer
	Sync() error
	Remove() error
}

// EmptyOutput is a writer that does nothing.
var EmptyOutput Writer = &emptyWriter{}

type emptyWriter struct{}

func (emptyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (emptyWriter) Remove() error {
	return nil
}

func (emptyWriter) Sync() error {
	return nil
}

// FileOutput return an file output.
func FileOutput(relativeDir, filename string) Writer {
	if !strings.Contains(filename, ".") {
		filename = fmt.Sprintf("%s.log", filename)
	}

	w, _ := os.OpenFile(fmt.Sprintf("%s/%s", getAbsPath(relativeDir), filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return &fileWriter{w}
}

type fileWriter struct {
	*os.File
}

func (w *fileWriter) Remove() error {
	if err := os.Remove(w.Name()); err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) {
			return nil
		}

		return err
	}

	return nil
}
