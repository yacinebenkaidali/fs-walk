package main

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
)

func filterOut(path string, minSize int64, ext string, info fs.FileInfo) bool {
	if !info.IsDir() && minSize > info.Size() {
		return true
	}
	if ext != "" && filepath.Ext(path) != ext {
		return true
	}
	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}
