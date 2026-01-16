package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
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

// to get the ReverseProxy configuration
type ProxyConfig struct {
	Port int `json:"port"`
	Strategy string `json:"strategy"` // e.g., "round-robin" or "least-conn"​
	HealthCheckFreq string `json:"health_check_frequency"`
}
type ReverseProxy struct{
	Port int
	Strategy string 
	HealthCheckFreq time.Duration 
	handler *httputil.ReverseProxy
}

type Config struct {
	ReverseProxy ProxyConfig `json:"proxy"`
	Backends []*Backend  `json: "backends"`
} 

func main() {
	configuration := getConfig("config.json")
	serverPool := initServerPool(configuration.Backends)
	

	//get backends list
	backends := serverPool.Backends
	
	//Reverse Proxy will forward to one of the servers in the ServerPool it's choosen randomly
	chosenIdx := rand.IntN(len(backends))
    targetURL := backends[chosenIdx].URL
	reverseProxy := newReverseProxy(configuration.ReverseProxy, targetURL)
	
	//Strating reverseProxy
	fmt.Println("Hello from reverseProxy, I'm starting now!")
	fmt.Printf("I'm forwarding to %s\n", backends[chosenIdx].URLString)

	http.Handle("/", reverseProxy.handler)
	http.ListenAndServe(":" + strconv.Itoa(reverseProxy.Port), nil)
}

// Functions used to get the configuration and initialize the Reverse Proxy and the server pool

func getConfig(fileName string)Config{
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	//unmarshall serverPool and Reverseproxy
	var config Config
	if err := json.Unmarshal(file, &config); err != nil{
		log.Fatal(err)
	}

	return config
}

func initServerPool(backends []*Backend) *ServerPool{
	serverPool := &ServerPool{
		Backends: backends,
	}
	for _, backend := range backends {
		parsedURL, err := url.Parse(backend.URLString)
		if err != nil{
			log.Fatal(err)
		}
		backend.URL = parsedURL
		backend.Alive = true
	}

	return serverPool
}

func newReverseProxy(proxyConfig ProxyConfig, targetURL *url.URL) *ReverseProxy{
	reverseProxy := &ReverseProxy{}
	reverseProxy.Port = proxyConfig.Port
	reverseProxy.Strategy = proxyConfig.Strategy
	reverseProxy.HealthCheckFreq, _ = time.ParseDuration(proxyConfig.HealthCheckFreq)
	reverseProxy.handler = httputil.NewSingleHostReverseProxy(targetURL)

	return reverseProxy
}
