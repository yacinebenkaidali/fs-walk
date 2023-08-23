package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	ext     string
	size    int64
	list    bool
	del     bool
	archive string
	log     io.Writer
}

var (
	f   = os.Stdout
	err error
)

func main() {
	root := flag.String("root", ".", "Root directory to start")
	list := flag.Bool("list", false, "List files only")
	ext := flag.String("ext", "", "File extension to filter out")
	logFile := flag.String("log", "", "Log deleted files to this file")
	size := flag.Int64("size", 0, "Minimum file size")
	del := flag.Bool("del", false, "Delete matched files")
	archiveDir := flag.String("archive", "", "Archive directory")

	flag.Parse()

	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	c := config{
		ext:     *ext,
		list:    *list,
		size:    *size,
		del:     *del,
		log:     f,
		archive: *archiveDir,
	}
	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, c config) error {
	delLogger := log.New(c.log, "Deleted file", log.LstdFlags)

	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filterOut(path, c.size, c.ext, info) {
			return nil
		}
		if c.list {
			return listFile(path, out)
		}
		if c.archive != "" {
			if err := archiveFile(c.archive, root, path); err != nil {
				return err
			}
		}
		if c.del {
			return delFile(path, delLogger)
		}
		return listFile(path, out)
	})
}
