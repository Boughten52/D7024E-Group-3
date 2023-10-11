package api

import (
	"d7024e/kademlia"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type API struct {
	kademlia *kademlia.Kademlia
}

// Create a new API instance.
func NewAPI(kademlia *kademlia.Kademlia) API {
	return API{kademlia}
}

// Handle POST request to upload objects.
func (api *API) UploadObjectHandler(w http.ResponseWriter, r *http.Request) {
	var content struct {
		Data string `json:"data"`
	}

	err := json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hash := api.kademlia.Store([]byte(content.Data))
	response := map[string]string{"hash": hash, "data": content.Data}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/objects/%s", hash)) // Set Location header
	w.WriteHeader(http.StatusCreated)                            // Set 201 Created status code
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

// Handle GET request to retrieve objects based on their hash.
func (api *API) GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	hash := strings.TrimPrefix(r.URL.Path, "/objects/")
	if len(hash) != 40 {
		http.Error(w, "Invalid hash length", http.StatusBadRequest)
		return
	}

	data := api.kademlia.LookupData(hash)
	if data == nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	response := map[string]string{"data": string(data)}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// Start the RESTful API server.
func StartServer(kademlia *kademlia.Kademlia, port int) {
	api := NewAPI(kademlia)
	http.HandleFunc("/objects", api.UploadObjectHandler) // Handle POST requests for uploading objects
	http.HandleFunc("/objects/", api.GetObjectHandler)   // Handle GET requests for retrieving objects by hash

	portStr := fmt.Sprintf("0.0.0.0:%d", port) // Listen on all interfaces
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
