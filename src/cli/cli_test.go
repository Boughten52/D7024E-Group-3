package cli

import (
	"testing"
)

func TestCLI_get(t *testing.T) {
	cli := NewCLI(nil, nil)

	testHash := "fake-hash"
	expectedOutput := ""

	capturedOutput := captureOutput(func() {
		cli.get(testHash)
	})

	if capturedOutput != expectedOutput {
		t.Errorf("get() did not produce the expected output when data was not found. Got: %s, Expected: %s", capturedOutput, expectedOutput)
	}
}

func captureOutput(f func()) string {
	oldOutput := output
	output = ""
	defer func() {
		output = oldOutput
	}()
	f()
	return output
}

var output string
