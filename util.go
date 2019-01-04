package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func OpenPuzzleInput(n int) io.ReadCloser {
	f, err := os.Open(filepath.Join(pkgdir(), "input", fmt.Sprintf("%02d.txt", n)))
	if err != nil {
		log.Fatalf("open puzzle input %02d: %v", n, err)
	}
	return f
}

func PuzzleInputLines(n int) []string {
	f := OpenPuzzleInput(n)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("reading puzzle input %02d: %v", n, err)
	}
	return lines
}

func pkgdir() string {
	_, fn, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}
	return filepath.Dir(fn)
}
