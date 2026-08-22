package sample

import (
	"fmt"
	alias "os"
	"path/filepath"
)

func Helper() int {
	return 1
}

type Worker struct {
	n int
}

func (w *Worker) Run() {
	_ = fmt.Sprint(w.n)
	_ = filepath.Base(".")
	_ = alias.Getenv("PATH")
}

type Counter interface {
	Inc()
}

type ID = string

func Main() {
	w := &Worker{}
	w.Run()
}
