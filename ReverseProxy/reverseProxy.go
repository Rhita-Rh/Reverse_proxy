package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"math/rand/v2"
)

type Backend struct{
	URL *url.URL `json:"-"`
	URLString string `json:"url"`
	Alive bool `json:"alive"`
	CurrentConns int64 `json:"current_connections"`
	mux sync.RWMutex
}

type ServerPool struct{
	Backends[ ]*Backend `json:"backends"`
	Current uint64 `json: "current"`
}

func main() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	//unmarshall serverPool
	var serverPool ServerPool
	err = json.Unmarshal(file, &serverPool)

	//get backends list
	backends := serverPool.Backends
	for _, backend := range serverPool.Backends {
		backend.URL, _ = url.Parse(backend.URLString)
		backend.Alive = true
	}
	//Reverse Proxy will forward to one of the servers in the ServerPool it's choosen randomly
	chosenIdx := rand.IntN(len(backends))
    targetURL := backends[chosenIdx].URL

	//Our reverse Proxy that may be customized after  
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	
	//Strating reverseProxy
	fmt.Println("Hello from reverseProxy, I'm starting now!")
	fmt.Printf("I'm forwarding to %s", backends[chosenIdx].URLString)

	http.Handle("/", reverseProxy)
	http.ListenAndServe(":8080", nil)
}
