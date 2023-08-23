package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "NoFilter", root: "testdata", cfg: config{ext: []string{""}, size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata", cfg: config{ext: []string{".log", ".sh"}, size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata", cfg: config{ext: []string{".log"}, size: 10, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata", cfg: config{ext: []string{".log"}, size: 20, list: true}, expected: ""},
		{name: "FilterExtensionNoMatch", root: "testdata", cfg: config{ext: []string{".gz"}, size: 0, list: true}, expected: ""},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var outBuffer bytes.Buffer
			err := run(tc.root, &outBuffer, tc.cfg)

			if err != nil {
				t.Fatal(err)
			}
			result := outBuffer.String()
			if result != tc.expected {
				t.Errorf("Expected %q and got %q instead\n", tc.expected, result)
			}
		})
	}
}
func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()
	dirname, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}
	for key, val := range files {
		for j := 0; j < val; j++ {
			fname := fmt.Sprintf("file%d%s", j, key)
			fpath := filepath.Join(dirname, fname)
			if err := os.WriteFile(fpath, []byte("dummy"), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}
	return dirname, func() {
		os.RemoveAll(dirname)
	}
}

func TestDelExtension(t *testing.T) {
	testCases := []struct {
		name string
		cfg  config
		//this is for all extensions
		nDelete     int
		nNoDelete   int
		expected    string
		extNoDelete string
	}{
		{name: "DeleteExtensionNoMatch", cfg: config{ext: []string{".log", ".sh"}, del: true}, extNoDelete: ".gz", nDelete: 0, nNoDelete: 10, expected: ""},
		{name: "DeleteExtensionMatch", cfg: config{ext: []string{".log"}, del: true}, extNoDelete: "", nDelete: 10, nNoDelete: 0, expected: ""},
		{name: "DeleteExtensionMixed", cfg: config{ext: []string{".log"}, del: true}, extNoDelete: ".gz", nDelete: 5, nNoDelete: 5, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			for _, ext := range tc.cfg.ext {

				var outBuffer bytes.Buffer
				var logBufer bytes.Buffer

				dirname, cleanup := createTempDir(t, map[string]int{
					ext:            tc.nDelete,
					tc.extNoDelete: tc.nNoDelete,
				})

				tc.cfg.log = &logBufer

				defer cleanup()
				if err := run(dirname, &outBuffer, tc.cfg); err != nil {
					t.Fatal(err)
				}
				res := outBuffer.String()

				if res != tc.expected {
					t.Errorf("Expected %q, got %q instead", tc.expected, res)
				}
				filesLeft, err := os.ReadDir(dirname)
				if err != nil {
					t.Fatal(err)
				}
				if len(filesLeft) != tc.nNoDelete {
					t.Errorf("Expected %d files left, got %d instead", tc.nNoDelete, len(filesLeft))
				}

				logLines := len(bytes.Split(logBufer.Bytes(), []byte("\n")))

				if logLines != tc.nDelete+1 {
					t.Errorf("Expected %d lines wrote to the log file, got %d instead", tc.nDelete, logLines)
				}
			}
		})
	}
}

func TestRunArchive(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{name: "ArchiveExtensionNoMatch", cfg: config{ext: []string{".log", ".sh"}}, extNoArchive: ".gz", nArchive: 0, nNoArchive: 10},
		{name: "ArchiveExtensionMatch", cfg: config{ext: []string{".log", ".someExt"}}, extNoArchive: "", nArchive: 10, nNoArchive: 0},
		{name: "ArchiveExtensionMixed", cfg: config{ext: []string{".log"}}, extNoArchive: ".gz", nArchive: 5, nNoArchive: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, ext := range tc.cfg.ext {

				var buffer bytes.Buffer
				tempDir, cleanup := createTempDir(t, map[string]int{
					ext:             tc.nArchive,
					tc.extNoArchive: tc.nNoArchive,
				})
				defer cleanup()
				archiveDir, cleanupArchive := createTempDir(t, nil)
				defer cleanupArchive()
				tc.cfg.archive = archiveDir

				if err := run(tempDir, &buffer, tc.cfg); err != nil {
					t.Fatal(err)
				}

				pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", ext))
				matchedFiles, err := filepath.Glob(pattern)
				if err != nil {
					t.Fatal(err)
				}

				expResult := strings.Join(matchedFiles, "\n")

				if strings.TrimSpace(buffer.String()) != expResult {
					t.Errorf("Expected %s, got %s instead", expResult, (&buffer).String())
				}
				files, err := os.ReadDir(archiveDir)
				if err != nil {
					t.Fatal(err)
				}
				if len(files) != tc.nArchive {
					t.Errorf("Expected %d files to be archived, got %d files archived instead", tc.nArchive, len(files))
				}
			}
		})
	}
}
