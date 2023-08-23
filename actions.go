package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"slices"
)

func filterOut(path string, minSize int64, ext []string, info fs.FileInfo) bool {
	if info.IsDir() || minSize > info.Size() {
		return true
	}
	//When no extension is selected return false to include all files
	if len(ext) == 1 && ext[0] == "" {
		return false
	}
	if !slices.Contains(ext, filepath.Ext(path)) {
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

func archiveFile(destDir, root, path string) error {
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	destPath := fmt.Sprintf("%s.gz", filepath.Base(path))

	targetPath := filepath.Join(destDir, relDir, destPath)

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}
	out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	zw := gzip.NewWriter(out)

	if _, err := io.Copy(zw, in); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}
	return out.Close()
}
