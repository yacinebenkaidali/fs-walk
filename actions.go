package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path string, minSize int64, ext string, info fs.FileInfo) bool {
	if info.IsDir() || minSize > info.Size() {
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

func delFile(path string, logger *log.Logger) error {

	if err := os.Remove(path); err != nil {
		return err
	}
	logger.Println(path)
	return nil
}
