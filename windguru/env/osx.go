package env

import (
	"log"
	"os"
	"strings"
)

func Getenv(key string) string {
	value := os.Getenv(key)
	if len(strings.TrimSpace(value)) == 0 {
		log.Panicf("no env key: %v", key)
	}

	return value
}
