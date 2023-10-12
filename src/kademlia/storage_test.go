package kademlia

import (
	"testing"
	"time"
)

func TestStorage_StoreData(t *testing.T) {
	storage := NewStorage(1 * time.Second) // Default TTL set to 1 second

	key := "test_key"
	data := []byte("test_data")

	// Test storing data for the first time
	storage.StoreData(key, data, 2*time.Second) // Custom TTL set to 2 seconds
	storedData, exists := storage.dataStore[key]
	if !exists {
		t.Errorf("Expected data to be stored for key %s, but it was not", key)
	}
	if string(storedData.Data) != "test_data" {
		t.Errorf("Expected data for key %s to be %s, but got %s", key, "test_data", string(storedData.Data))
	}

	// Test storing data for the same key again with a shorter TTL
	storage.StoreData(key, []byte("new_test_data"), 1*time.Second) // Custom TTL set to 1 second
	storedData, exists = storage.dataStore[key]
	if !exists {
		t.Errorf("Expected data to be stored for key %s, but it was not", key)
	}
	if string(storedData.Data) != "test_data" {
		t.Errorf("Expected data for key %s to remain unchanged, but it was modified", key)
	}
}

func TestStorage_FetchData(t *testing.T) {
	storage := NewStorage(1 * time.Second) // Default TTL set to 1 second

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
	storage.dataStore[key] = struct {
		Data []byte
		TTL  time.Time
	}{Data: data, TTL: time.Now().Add(1 * time.Second)} // Data with TTL set to 1 second
	result, exists = storage.FetchData(key)
	if !exists {
		t.Errorf("Expected data for key %s to exist, but it was not found", key)
	}
	if string(result) != string(data) {
		t.Errorf("Expected fetched data to match the stored data, but it did not")
	}
}

func TestStorage_RefreshDataTTL(t *testing.T) {
	storage := &Storage{
		dataStore: make(map[string]struct {
			Data []byte
			TTL  time.Time
		}),
	}

	// Add a data object to the storage with a specific TTL
	key := "test_key"
	data := []byte("test_data")
	ttl := 2 * time.Second
	expirationTime := time.Now().Add(ttl)
	storage.dataStore[key] = struct {
		Data []byte
		TTL  time.Time
	}{Data: data, TTL: expirationTime}

	// Attempt to refresh TTL for the existing data object
	newTTL := 4 * time.Second
	refreshed := storage.RefreshDataTTL(key, newTTL)

	// Check if TTL was successfully refreshed
	if !refreshed {
		t.Errorf("Expected TTL to be refreshed, but it was not")
	}

	// Check if TTL was updated in the data store
	storedData, exists := storage.dataStore[key]
	if !exists {
		t.Errorf("Expected data for key %s to exist, but it was not found", key)
	}

	// Calculate the expected expiration time after refresh
	expectedExpirationTime := time.Now().Add(newTTL)

	// Allow for a small time difference (1 millisecond) due to execution time variations
	timeDifference := storedData.TTL.Sub(expectedExpirationTime)
	if timeDifference > time.Millisecond {
		t.Errorf("Expected TTL to be updated, but it remains unchanged")
	}

	// Attempt to refresh TTL for a non-existing data object
	nonExistingKey := "non_existing_key"
	refreshed = storage.RefreshDataTTL(nonExistingKey, newTTL)

	// Check if TTL was not refreshed for a non-existing data object
	if refreshed {
		t.Errorf("Expected TTL not to be refreshed for a non-existing key, but it was refreshed")
	}
}
