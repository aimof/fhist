package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		return
	}
	fHistDir := os.Getenv("FHISTDIR")

	fHistAbsDir, err := filepath.Abs(fHistDir)
	if err != nil {
		fHistAbsDir = os.Getenv("HOME") + "/.fhist"
		return
	}
	info, err := os.Stat(fHistAbsDir)
	if err != nil {
		return
	} else if !info.IsDir() {
		return
	}

	switch args[1] {
	case "save":
		save(args, fHistAbsDir)
	}
}

func save(args []string, fHistAbsDir string) {
	if len(args) < 3 {
		return
	}

	splitBuffer := strings.Split(args[2], " ")
	cmd := make([]string, 0)
	for _, s := range splitBuffer {
		if s != "" {
			cmd = append(cmd, s)
		}
	}
	if len(cmd) == 1 {
		return
	}
	file, err := os.OpenFile(path.Clean(fHistAbsDir+"/"+cmd[0]), os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return
	}
	for i := 1; i < len(cmd); i++ {
		abs, err := filepath.Abs(cmd[i])
		if err != nil {
			continue
		}
		info, err := os.Stat(abs)
		if info.IsDir() {
			if !(abs[len(abs)-1:] == "/") {
				abs = abs + "/"
			}
		}
		fmt.Fprintln(file, abs)
	}
	return
}
