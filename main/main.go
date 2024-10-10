package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	parser "github.com/Re1nGer/go_jp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run script.go <path_to_json_files>")
		os.Exit(1)
	}

	path := os.Args[1]

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			filename := file.Name()
			filePath := filepath.Join(path, filename)
			test_should_fail := strings.Contains(filename, "fail")
			fmt.Printf("Processing file: %s\n", filePath)

			// Read file content
			content, err := os.ReadFile(filePath)
			if err != nil && !test_should_fail {
				fmt.Printf("Error reading file %s: %v\n", filePath, err)
				continue
			}

			fmt.Println("content", string(content))

			p, err := parser.NewParser(content)

			if err != nil {
				fmt.Printf("Error parsing JSON in file %s: %v\n %v\n", filePath, err)
				continue
			}

			_, err = p.Parse()

			if err != nil {
				fmt.Printf("Error parsing JSON in file %s: %v\n %v", filePath, err)
			} else {
				fmt.Printf("Successfully parsed JSON in file %s\n %v\n", filePath)
			}
		}
	}
}
