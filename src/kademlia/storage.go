package kademlia

import (
	"d7024e/utils"
	"sync"
	"time"
)

type Storage struct {
	mu        sync.Mutex
	dataStore map[string]struct {
		Data []byte
		TTL  time.Time
	}
	DefaultTTL time.Duration
}

// Initializes the Storage struct with a default TTL value
func NewStorage(defaultTTL time.Duration) *Storage {
	storage := &Storage{
		dataStore: make(map[string]struct {
			Data []byte
			TTL  time.Time
		}),
		DefaultTTL: defaultTTL,
	}

	// Start a goroutine to periodically clean up expired objects
	go storage.startCleanupTask()
	return storage
}

// Stores data locally but does not overwrite any already defined key data pairs
func (storage *Storage) StoreData(key string, data []byte, ttl time.Duration) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	existingData, exist := storage.dataStore[key]
	if exist && time.Now().Before(existingData.TTL) {
		utils.Log(3, "Tried to store data %s with key %s but that key already stores data %s", string(data), key, existingData.Data)
		return
	}

	expirationTime := time.Now().Add(ttl)
	storage.dataStore[key] = struct {
		Data []byte
		TTL  time.Time
	}{Data: data, TTL: expirationTime}

	utils.Log(1, "Successfully stored data %s with key %s (TTL: %s)", string(data), key, expirationTime.String())
}

// Tries to retrieve data and returns it together with the success of the fetch
func (storage *Storage) FetchData(key string) ([]byte, bool) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	existingData, exist := storage.dataStore[key]
	if exist && time.Now().Before(existingData.TTL) {
		// Reset TTL since data object is requested
		storage.dataStore[key] = struct {
			Data []byte
			TTL  time.Time
		}{Data: existingData.Data, TTL: time.Now().Add(storage.DefaultTTL)}
		return existingData.Data, true
	}

	// Delete the data object if TTL has expired
	delete(storage.dataStore, key)
	return nil, false
}

// Refreshes the TTL for a data object if it exists and has not expired. Returns true if the TTL was refreshed.
func (storage *Storage) RefreshDataTTL(key string, ttl time.Duration) bool {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if storedData, exists := storage.dataStore[key]; exists && time.Now().Before(storedData.TTL) {
		// Reset TTL for the data object
		storage.dataStore[key] = struct {
			Data []byte
			TTL  time.Time
		}{Data: storedData.Data, TTL: time.Now().Add(ttl)}
		return true
	}
	return false
}

// Periodically checks and deletes expired objects from the data store
func (storage *Storage) startCleanupTask() {
	ticker := time.NewTicker(storage.DefaultTTL)
	for range ticker.C {
		storage.mu.Lock()
		for key, data := range storage.dataStore {
			if time.Now().After(data.TTL) {
				utils.Log(2, "Deleting expired data with key %s", key)
				delete(storage.dataStore, key)
			}
		}
		storage.mu.Unlock()
	}
}
