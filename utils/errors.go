package utils

import (
	"fmt"
	"os"
)

func Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message)
	os.Exit(1)
}
