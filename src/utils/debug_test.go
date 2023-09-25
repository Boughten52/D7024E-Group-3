package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func captureOutput(f func()) string {
	oldOutput := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)

	// Copy the output in a separate goroutine so it can be captured
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	f()

	// Restore the original stdout
	os.Stdout = oldOutput
	w.Close()

	return <-outC
}

func TestLog(t *testing.T) {
	tests := []struct {
		severity   int
		threshold  int
		message    string
		expected   string
		shouldLog  bool
		shouldSkip bool
	}{
		{1, 1, "Test Message", "Test Message\n", true, false},
		{5, 1, "Test Message", "Test Message\n", true, false},
		{10, 1, "Test Message", "Test Message\n", true, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Severity %d", test.severity), func(t *testing.T) {
			severityThreshold = test.threshold
			output := captureOutput(func() {
				Log(test.severity, test.message)
			})

			if (output == test.expected) != test.shouldLog {
				t.Errorf("Log() for severity %d did not produce the expected output", test.severity)
			}

			if output != test.expected {
				t.Errorf("Log() for severity %d produced: %s, expected: %s", test.severity, output, test.expected)
			}
		})
	}
}

func TestLogError(t *testing.T) {
	tests := []struct {
		severity   int
		threshold  int
		message    string
		expected   string
		shouldLog  bool
		shouldSkip bool
	}{
		{1, 1, "Error Message", "ERROR: Error Message\n", true, false},
		{5, 1, "Error Message", "ERROR: Error Message\n", true, false},
		{10, 1, "Error Message", "ERROR: Error Message\n", true, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Severity %d", test.severity), func(t *testing.T) {
			severityThreshold = test.threshold
			output := captureOutput(func() {
				LogError(test.message)
			})

			if (output == test.expected) != test.shouldLog {
				t.Errorf("LogError() for severity %d did not produce the expected output", test.severity)
			}

			if output != test.expected {
				t.Errorf("LogError() for severity %d produced: %s, expected: %s", test.severity, output, test.expected)
			}
		})
	}
}
