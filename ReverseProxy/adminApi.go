package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type Status struct{
	Total_backends int `json:"total_backends"`
	Active_backends int `json:"active_backends"`
	Backends []BackendResponse `json:"backends"`
}

type BackendResponse struct{
	URLString string `json:"url"`
	Alive bool `json:"alive"`
	CurrentConns int64 `json:"current_connections"`
}

func validURL(URL *url.URL) bool{
	return URL.Host != "" && URL.Scheme != ""
}

func statusHandleGet(w http.ResponseWriter, _ *http.Request, serverPool *ServerPool){
	serverPool.mux.Lock() 	//locking serverPool to read it for status
	defer serverPool.mux.Unlock()

	var status Status
	status.Total_backends = len(serverPool.Backends)

	for _, backend := range serverPool.Backends{
		alive := backend.IsAlive()
		if alive{
			status.Active_backends ++
		}

		backend_response := BackendResponse{
			URLString: backend.URLString,
			Alive: alive,
			CurrentConns: backend.GetCurentConn(),
		}
		status.Backends = append(status.Backends, backend_response)
	}

	jsonResp, err := json.Marshal(status)
	if err != nil {
		http.Error(w, "Failed to marshal status", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func updateConfigFile(w http.ResponseWriter, r *http.Request, backendPost BackendConfig, indexToDelete int, fileName string){
	// update configuration file
	configFile, err := os.ReadFile(fileName)
    if err != nil {
        http.Error(w, "Failed to read config file", http.StatusInternalServerError)
        return
    }
	var config Config
    if err := json.Unmarshal(configFile, &config); err != nil {
        http.Error(w, "Failed to parse config file", http.StatusInternalServerError)
        return
    }

	if r.Method == http.MethodPost{
		config.Backends = append(config.Backends, &backendPost)
	}else if r.Method == http.MethodDelete{
		config.Backends = append(config.Backends[:indexToDelete], config.Backends[indexToDelete+1:] ...)
	}
	

	// Marshal to JSON
	updatedBackends, err := json.MarshalIndent(config, "", " ")
	if err!=nil{
		http.Error(w, "Failed to add to JSON", http.StatusInternalServerError)
        return
	}

	 // Write to file
    if err := os.WriteFile("config.json", updatedBackends, 0644); err != nil {
        http.Error(w, "Failed to write config file", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

func backendHandlePost(w http.ResponseWriter, r *http.Request, serverPool *ServerPool){
	if r.Method != http.MethodPost{
		http.Error(w, "Error", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	// Decode the posted backend
	var backendPost BackendConfig
	err := json.NewDecoder(r.Body).Decode(&backendPost)
	if err != nil{
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Verify posted URL (is it a valid URL)
	parsedURL, err := url.Parse(backendPost.URLString)
	if err!= nil || !validURL(parsedURL){
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	
	// Lock serverPool to add the backend
	serverPool.mux.Lock()
	defer serverPool.mux.Unlock()

	// Create the Backend with the corresponding URL
	backendToAdd := Backend{
		URLString: backendPost.URLString,
		URL: parsedURL,
		Alive: false,
		CurrentConns: 0,
	}

	// prevent duplicates
	for _, backend := range serverPool.Backends{
		if backend.URLString == backendToAdd.URLString{
			http.Error(w, "URL duplicated", http.StatusBadRequest)
			return
		}
	}

	serverPool.Backends = append(serverPool.Backends, &backendToAdd)

	// update configuration file
	updateConfigFile(w, r, backendPost, 0, "config.json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Backend added successfully",
		"url": backendToAdd.URLString,
	})
}

func backendHandleDelete(w http.ResponseWriter, r *http.Request, serverPool *ServerPool){
	if r.Method != http.MethodDelete{
		http.Error(w, "Error", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var backendToDelete BackendConfig
	err := json.NewDecoder(r.Body).Decode(&backendToDelete)
	if err != nil{
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Verify posted URL (is it a valid URL)
	parsedURL, err := url.Parse(backendToDelete.URLString)
	if err!= nil || !validURL(parsedURL){
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	serverPool.mux.Lock()
	defer serverPool.mux.Unlock()

	idxToDelete := len(serverPool.Backends)
	for i, backend := range serverPool.Backends{
		if backend.URLString == backendToDelete.URLString{
			idxToDelete = i
			break
		}
	}
	if idxToDelete == len(serverPool.Backends){
		http.Error(w, "non-existing Backend", http.StatusBadRequest)
		return
	}
	serverPool.Backends = append(serverPool.Backends[:idxToDelete], serverPool.Backends[idxToDelete+1:] ...)
	
	// update configuration file
	updateConfigFile(w, r, BackendConfig{}, idxToDelete, "config.json")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backend deleted successfully\n"))
}

func backendHandler(w http.ResponseWriter, r *http.Request, serverPool *ServerPool){
	switch r.Method{
	case http.MethodPost:
		backendHandlePost(w, r, serverPool)
	case http.MethodDelete:
		backendHandleDelete(w, r, serverPool)
	default:
		http.Error(w, "Method not yet implemented", http.StatusMethodNotAllowed)
	}
}

func AdminApi(serverPool *ServerPool){
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request){
		if r.Method != http.MethodGet{
			http.Error(w, "Error", http.StatusMethodNotAllowed)
			return
		}
		statusHandleGet(w, r, serverPool)
	})

	serverMux.HandleFunc("/backends", func(w http.ResponseWriter, r *http.Request){
		backendHandler(w, r, serverPool)
	})

	fmt.Println("Hello from Admin API. I'm listening on port :9000")
	http.ListenAndServe(":9000", serverMux)
}