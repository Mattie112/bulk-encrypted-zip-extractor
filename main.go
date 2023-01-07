package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var path = ""
var BinPath = ""
var passwordList = []string{""}
var del = false

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) < 3 {
		fmt.Println("Use <dir_of_files> <path_to_7zz> <path_to_passwords.txt> <optional: delete>")
	}
	path = argsWithoutProg[0]
	BinPath = argsWithoutProg[1]
	passwordPath := argsWithoutProg[2]
	delStr := argsWithoutProg[3]
	if delStr == "true" {
		del = true
	}
	checkBinary()
	readPasswordsFromFile(passwordPath)
	files, _ := os.ReadDir(path)
	for _, file := range files {
		for _, ext := range getSupportedExtensions() {
			if strings.Contains(file.Name(), ext) {
				fmt.Println(file.Name(), file.Type().String())
				extractFile(path + "/" + file.Name())
			}
		}
	}
	fmt.Println("Done")
}

func readPasswordsFromFile(filePath string) {
	readFile, err := os.Open(filePath)

	if err != nil {
		fmt.Println(err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		password := fileScanner.Text()
		fmt.Println(fmt.Sprintf("Found password: '%s'", password))
		passwordList = append(passwordList, password)
	}

	_ = readFile.Close()
}

func extractFile(filePath string) {
	fmt.Println("Starting the extraction of: " + filePath)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	for _, password := range passwordList {
		fmt.Println(fmt.Sprintf("Trying %s with '%s'", filePath, password))
		args := []string{"e", fmt.Sprintf("-p%s", password), fmt.Sprintf("-o%s", path), "-y", filePath}
		cmd := exec.CommandContext(ctx, BinPath, args...)
		err := cmd.Run()
		if err == nil {
			if del {
				fmt.Println("Extract OK, removing file")
				err := os.Remove(filePath)
				if err != nil {
					fmt.Println(fmt.Sprintf("Could not delete %s, will continue", filePath))
				}
			}
			break
		}
		if exiterr, ok := err.(*exec.ExitError); ok {
			log.Printf("Exit Status: %d", exiterr.ExitCode())
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	fmt.Println("Extract complete")
}

func checkBinary() {
	if BinPath != "" {
		// Check if we already have a full path to the executeable
		if _, err := os.Stat(BinPath); err == nil {
			return
		}
	}

	binPath, err := exec.LookPath("7zz")
	if err == nil {
		BinPath = binPath
		return
	}
	if _, err := os.Stat(filepath.Join(path, "7zz")); err == nil {
		BinPath = filepath.Join(path, "7zz")
		return
	}
	panic("Could not find 7zz binary")
}

func getSupportedExtensions() []string {
	return []string{".zip", ".rar", ".7z"}
}
