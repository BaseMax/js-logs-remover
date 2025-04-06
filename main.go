package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"regexp"
	"sync"

	"github.com/fatih/color"
)

var consoleMethods = []string{
	"assert", "clear", "count", "countReset", "debug", "dir", "dirxml",
	"error", "group", "groupCollapsed", "groupEnd", "info", "log",
	"table", "time", "timeEnd", "timeLog", "timeStamp", "trace", "warn",
}

var methodRegex *regexp.Regexp
var allowedMethods map[string]bool

func init() {
	methodRegex = regexp.MustCompile(`console\.(\w+)`)
	allowedMethods = make(map[string]bool)
}

func removeConsoleLogsFromFile(filePath string) error {
	code, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(code)
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		for _, method := range allowedMethods {
			if strings.Contains(line, "console."+method+"(") {
				lines[i] = ""
				break
			}
		}
	}
	updatedCode := strings.Join(lines, "\n")
	return ioutil.WriteFile(filePath, []byte(updatedCode), 0644)
}

func processFile(filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	ext := filepath.Ext(filePath)
	if ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx" {
		err := removeConsoleLogsFromFile(filePath)
		if err != nil {
			color.Red("Error processing file %s: %v", filePath, err)
		} else {
			color.Green("Processed: %s", filePath)
		}
	}
}

func processDirectory(dirPath string, wg *sync.WaitGroup) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		color.Red("Error reading directory %s: %v", dirPath, err)
		return
	}

	for _, file := range files {
		fullPath := filepath.Join(dirPath, file.Name())
		if file.IsDir() {
			if file.Name() != "node_modules" && file.Name() != ".git" && file.Name() != "dist" && file.Name() != "build" {
				processDirectory(fullPath, wg)
			}
		} else {
			wg.Add(1)
			go processFile(fullPath, wg)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: js-logs-remover [path] [log-methods]")
		return
	}

	args := os.Args[2:]
	removeAllLogs := false
	if len(args) == 1 && args[0] == "all" {
		removeAllLogs = true
	}

	if removeAllLogs {
		for _, method := range consoleMethods {
			allowedMethods[method] = true
		}
	} else {
		for _, arg := range args {
			methods := strings.Split(arg, ",")
			for _, method := range methods {
				allowedMethods[method] = true
			}
		}
	}

	targetDir := os.Args[1]
	if targetDir == "" {
		targetDir = "."
	}

	var wg sync.WaitGroup
	processDirectory(targetDir, &wg)
	wg.Wait()

	color.Green("✅ All selected console methods removed!")
}
