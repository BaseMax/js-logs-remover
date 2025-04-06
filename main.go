package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"
	"sync"

	"github.com/fatih/color"
)

var consoleMethods = []string{
	"assert", "clear", "count", "countReset", "debug", "dir", "dirxml",
	"error", "group", "groupCollapsed", "groupEnd", "info", "log",
	"table", "time", "timeEnd", "timeLog", "timeStamp", "trace", "warn",
}

var allowedMethods map[string]bool

func init() {
	allowedMethods = make(map[string]bool)
}

func removeConsoleLogsFromFile(filePath string) error {
	code, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	content := string(code)
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		for method := range allowedMethods {
			if strings.Contains(line, "console."+method+"(") {
				lines[i] = ""
				break
			}
		}
	}

	updatedCode := strings.Join(lines, "\n")
	if err := ioutil.WriteFile(filePath, []byte(updatedCode), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return nil
}

func processFile(filePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	ext := filepath.Ext(filePath)
	if ext == ".js" || ext == ".ts" || ext == ".jsx" || ext == ".tsx" {
		if err := removeConsoleLogsFromFile(filePath); err != nil {
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
			if !isExcludedDir(file.Name()) {
				processDirectory(fullPath, wg)
			}
		} else {
			wg.Add(1)
			go processFile(fullPath, wg)
		}
	}
}

func isExcludedDir(dirName string) bool {
	excludedDirs := map[string]bool{
		"node_modules": true,
		".git":          true,
		"dist":          true,
		"build":         true,
	}
	return excludedDirs[dirName]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: js-logs-remover [path] [log-methods]")
		return
	}

	args := os.Args[2:]
	if len(args) == 1 && args[0] == "all" {
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
