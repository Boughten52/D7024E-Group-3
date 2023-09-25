package kademlia

import (
	"testing"
)

func TestStorage_StoreData(t *testing.T) {
	storage := &Storage{
		dataStore: make(map[string][]byte),
	}

	key := "test_key"
	data := []byte("test_data")

	// Test storing data for the first time
	storage.StoreData(key, data)
	if _, ok := storage.dataStore[key]; !ok {
		t.Errorf("Expected data to be stored for key %s, but it was not", key)
	}

	// Test storing data for the same key again
	storage.StoreData(key, []byte("new_test_data"))
	if string(storage.dataStore[key]) != "test_data" {
		t.Errorf("Expected data for key %s to remain unchanged, but it was modified", key)
	}
}

func TestStorage_FetchData(t *testing.T) {
	storage := &Storage{
		dataStore: make(map[string][]byte),
	}

	key := "test_key"
	data := []byte("test_data")

	// Test fetching data that does not exist
	result, exists := storage.FetchData(key)
	if exists {
		t.Errorf("Expected data for key %s not to exist, but it was found", key)
	}
	if result != nil {
		t.Errorf("Expected fetched data to be nil, but it was not")
	}

	// Test fetching data that exists
	storage.dataStore[key] = data
	result, exists = storage.FetchData(key)
	if !exists {
		t.Errorf("Expected data for key %s to exist, but it was not found", key)
	}
	if string(result) != string(data) {
		t.Errorf("Expected fetched data to match the stored data, but it did not")
	}
}
