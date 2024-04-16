package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("[-] Open %s error, %v\n.", filename, err)
		os.Exit(0)
	}
	defer file.Close()

	var content []string
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text != "" {
			content = append(content, scanner.Text())
		}
	}

	return content, nil
}
