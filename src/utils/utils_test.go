package utils

import (
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
	expectedHash := "2ef7bde608ce5404e97d5f042f95f89f1c1e6b962"

	hash := Hash(data)
	if hash != expectedHash {
		t.Errorf("Hash() returned %s, expected %s", hash, expectedHash)
	}
}
