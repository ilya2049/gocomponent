package fs

import (
	"fmt"
	"os"
	"strings"
)

func isGoSourceFile(path string) bool {
	return strings.Contains(path, ".go")
}

func readFile(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read: %w", err)
	}

	return string(contents), nil
}
