package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// RedactingWriter wraps an underlying io.Writer (e.g., os.Stdout)
type RedactingWriter struct {
	Destination io.Writer
}

// Implement the io.Writer interface on RedactingWriter
func (rw *RedactingWriter) Write(p []byte) (int, error) {
	// 1. Replace sensitive keyword "ERROR" with "[REDACTED_ERROR]"
	sanitizedBytes := bytes.ReplaceAll(p, []byte("ERROR"), []byte("[REDACTED_ERROR]"))

	// 2. Pass the sanitized bytes down to the target destination
	_, err := rw.Destination.Write(sanitizedBytes)
	if err != nil {
		return 0, err
	}

	// 3. Return original len(p) so callers know all input bytes were processed
	return len(p), nil
}

func main() {
	// Open source log
	file, err := os.Open("sample.log")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Wrap os.Stdout inside our custom RedactingWriter
	redactor := &RedactingWriter{Destination: os.Stdout}

	fmt.Println("--- Streaming sample.log through Custom RedactingWriter ---")

	// Stream file directly into our redactor!
	_, err = io.Copy(redactor, file)
	if err != nil {
		fmt.Printf("Error during copy: %v\n", err)
		return
	}
}
