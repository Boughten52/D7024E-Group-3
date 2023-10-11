package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"
)

func TestGetIP(t *testing.T) {
	ip, err := GetIP()
	if err != nil {
		t.Errorf("GetIP() returned an error: %v", err)
	}

	if ip == "" {
		t.Error("GetIP() returned an empty IP address")
	}
}

func TestHash(t *testing.T) {
	data := []byte("hello world")

	// Calculate expected hash using the standard library SHA-1 implementation
	hasher := sha1.New()
	hasher.Write(data)
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	// Calculate hash using your Hash() function
	hash := Hash(data)

	// Test consistency
	if hash != expectedHash {
		t.Errorf("Hash() returned %s, expected %s", hash, expectedHash)
	}

	// Test avalanche effect
	data[0] = 'H' // Change the first byte of input
	hasher = sha1.New()
	hasher.Write(data)
	newExpectedHash := hex.EncodeToString(hasher.Sum(nil))

	if hash == newExpectedHash {
		t.Errorf("Hash() did not change with a small input change")
	}
}
