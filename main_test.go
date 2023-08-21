package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{name: "NoFilter", root: "testdata", cfg: config{ext: "", size: 0, list: true}, expected: "testdata/dir.log\ntestdata/dir2/script.sh\n"},
		{name: "FilterExtensionMatch", root: "testdata", cfg: config{ext: ".log", size: 0, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeMatch", root: "testdata", cfg: config{ext: ".log", size: 10, list: true}, expected: "testdata/dir.log\n"},
		{name: "FilterExtensionSizeNoMatch", root: "testdata", cfg: config{ext: ".log", size: 20, list: true}, expected: ""},
		{name: "FilterExtensionNoMatch", root: "testdata", cfg: config{ext: ".gz", size: 0, list: true}, expected: ""},
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
		name        string
		cfg         config
		nDelete     int
		nNoDelete   int
		expected    string
		extNoDelete string
	}{
		{name: "DeleteExtensionNoMatch", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 0, nNoDelete: 10, expected: ""},
		{name: "DeleteExtensionMatch", cfg: config{ext: ".log", del: true}, extNoDelete: "", nDelete: 10, nNoDelete: 0, expected: ""},
		{name: "DeleteExtensionMixed", cfg: config{ext: ".log", del: true}, extNoDelete: ".gz", nDelete: 5, nNoDelete: 5, expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var outBuffer bytes.Buffer
			dirname, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})

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
		})
	}
}
