package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// ReadLines eagerly reads all lines from filePath.
func ReadLines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %w", filePath, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func GetModuleRootPath() string {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}

	output := string(outputBytes)
	if output[len(output)-2] == '\r' {
		return output[0 : len(output)-2]
	}
	return output[0 : len(output)-1]
}
