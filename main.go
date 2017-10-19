package main

import (
	"bufio"
	"fmt"
	"log"
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
	if fHistDir == "" {
		fHistDir = os.Getenv("HOME") + "/.fhist"
	}
	fHistAbsDir, err := filepath.Abs(fHistDir)
	if err != nil {
		log.Println("cannot find $FHISTDIR")
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
		err = save(args, fHistAbsDir)
		if err != nil {
			log.Fatalln(err)
		}
	case "list":
		err = list(args, fHistAbsDir)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func save(args []string, fHistAbsDir string) error {
	if len(args) < 3 {
		return nil
	}

	cmd := parseCmd(args[2])
	if len(cmd) == 1 {
		return nil
	}
	file, err := os.OpenFile(path.Clean(fHistAbsDir+"/"+cmd[0]), os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	for i := 1; i < len(cmd); i++ {
		abs, err := filepath.Abs(cmd[i])
		if err != nil {
			continue
		}
		info, err := os.Stat(abs)
		if err != nil {
			continue
		}
		if info.IsDir() {
			if !(abs[len(abs)-1:] == "/") {
				abs = abs + "/"
			}
		}
		fmt.Fprintln(file, abs)
	}
	return nil
}

func list(args []string, fHistAbsDir string) error {
	if len(args) < 3 {
		return nil
	}
	cmd := parseCmd(args[2])
	if len(cmd) == 0 {
		return nil
	}

	file, err := os.OpenFile(path.Clean(fHistAbsDir+"/"+cmd[0]), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	absPaths := make([]string, 0)
	relPaths := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		absTarget, err := filepath.Abs(scanner.Text())
		if err != nil {
			return err
		}
		if !strings.Contains(absTarget, wd) {
			absPaths = append(absPaths, absTarget)
			continue
		}
		relPath, err := filepath.Rel(wd, scanner.Text())
		if err != nil {
			continue
		}
		if relPath != "." {
			relPath = "./" + relPath
		}
		relPaths = append(relPaths, relPath)
	}
	var output string
	for _, p := range relPaths {
		output = output + p + "\n"
	}
	for _, p := range absPaths {
		output = output + p + "\n"
	}
	fmt.Print(output)
	return nil
}

func parseCmd(rawCmd string) []string {
	splitBuffer := strings.Split(rawCmd, " ")
	cmd := make([]string, 0)
	for _, s := range splitBuffer {
		if s != "" {
			cmd = append(cmd, s)
		}
	}
	return cmd
}
