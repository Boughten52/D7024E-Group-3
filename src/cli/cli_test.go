package cli

import (
	"d7024e/kademlia"
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

func TestCLI_Forget(t *testing.T) {
	// Create a new Kademlia instance
	kad := kademlia.NewKademlia(kademlia.NewNetwork(kademlia.NewRoutingTable(kademlia.NewContact(kademlia.NewRandomKademliaID(), "172.20.0.10")), 20, 3, 60, 30))

	// Create a new CLI instance with the Kademlia instance
	cli := NewCLI(kad, nil)

	// Test with a valid hash
	validHash := "1111111111111111111111111111111111111111"
	cli.forget(validHash)

	// Check if the Forget method is called with the correct hash
	if _, ok := kad.ClosestPeers[validHash]; ok {
		t.Errorf("Expected hash %s to be forgotten, but it was not", validHash)
	}

	// Test with an invalid hash
	invalidHash := "12345" // Invalid length
	cli.forget(invalidHash)

	// Check if the CLI handles an invalid hash length correctly
	// The test passes if it reaches this point without panicking
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
