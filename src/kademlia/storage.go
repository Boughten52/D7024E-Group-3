package kademlia

import "d7024e/utils"

type Storage struct {
	dataStore map[string][]byte
}

// Stores data locally but does not overwrite any already defined key data pairs
func (storage *Storage) StoreData(key string, data []byte) {

	_, exist := storage.dataStore[key]
	if exist {
		utils.Log(3, "Tried to store data %s with key %s but that key already stores data %s", string(data), key, storage.dataStore[key])
		return
	}

	utils.Log(1, "Successfully stored data %s with key %s", string(data), key)
	storage.dataStore[key] = data
}

// Tries to retrieve data and returns it together with the success of the fetch
func (storage *Storage) FetchData(key string) ([]byte, bool) {

	_, exist := storage.dataStore[key]
	if exist {
		return storage.dataStore[key], true
	}

	return nil, false
}
