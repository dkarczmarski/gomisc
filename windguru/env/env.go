// Package env provides methods to read environment variables
package env

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

func Load(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("env file open error: %w", err)
	}

	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		envKey := parts[0]
		envValue := parts[1]

		if len(os.Getenv(envKey)) > 0 {
			continue
		}

		if err := os.Setenv(envKey, envValue); err != nil {
			return err
		}
	}

	return nil
}
