package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

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
	ServerPool *ServerPool
}

type Config struct {
	ReverseProxy ProxyConfig `json:"proxy"`
	Backends []*Backend  `json:"backends"`
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

func newReverseProxy(proxyConfig ProxyConfig, serverPool *ServerPool) *ReverseProxy{
	reverseProxy := &ReverseProxy{}
	reverseProxy.Port = proxyConfig.Port
	reverseProxy.Strategy = proxyConfig.Strategy
	reverseProxy.HealthCheckFreq, _ = time.ParseDuration(proxyConfig.HealthCheckFreq)
	reverseProxy.ServerPool = serverPool

	return reverseProxy
}

// ReverseProxy handler
func (reverseProxy *ReverseProxy)ServeHTTP(w http.ResponseWriter, r *http.Request){
	// request context is passed through to cancel slow backend processing
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	r = r.WithContext(ctx)

	var validBackend *Backend
	switch reverseProxy.Strategy{
	case "round-robin":
		validBackend = reverseProxy.ServerPool.GetNextValidPeer_RoundRobin()
	case"least-conn":
		validBackend = reverseProxy.ServerPool.GetNextValidPeer_LeastConn()
	default:
		http.Error(w, "Strategy not yet implemented: choose round-robin or least-conn", 500)
		return
	}
	if validBackend == nil{
		http.Error(w, "Service Unavailable", 503)
		return
	}

	goProxy := httputil.NewSingleHostReverseProxy(validBackend.URL)
	fmt.Printf("I'm forwarding to %s\n", validBackend.URLString)
	//increment number of connections of the server when request comes
	validBackend.IncrementConn()

	//decrement whenever request is handeled 
	defer validBackend.DecrementConn()

	//forward request to server
	goProxy.ServeHTTP(w, r)

}

// Health-checking function
func HealthCheckFunc(reverseProxy *ReverseProxy, periodicity time.Duration){
	serverPool := reverseProxy.ServerPool
	for{
		// Health for backends concurrently
		var wg sync.WaitGroup
		fmt.Println("Health Check Starts now!")
		for _, backend := range serverPool.Backends{
			wg.Add(1)
			go func(){
				defer wg.Done()
				address := backend.URL.Host
				alive := pingTCP(address)
				serverPool.SetBackendStatus(backend.URL, alive)
				if alive{
					fmt.Printf("Backend %s is UP\n", backend.URLString)
				}else{
					fmt.Printf("Backend %s is DOWN\n", backend.URLString)
				}
			}()
		}
		wg.Wait()
		time.Sleep(periodicity)
	}	
}

// ping using TCP dial
func pingTCP(address string) bool{
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil{
		return false
	}

	conn.Close()
	return true
	
}

//==============================================================================

// Main func
func main() {
	configuration := getConfig("config.json")
	serverPool := initServerPool(configuration.Backends)
	
	//Reverse Proxy will pass the request to the loadbalancer
	reverseProxy := newReverseProxy(configuration.ReverseProxy, serverPool)

	//Starting reverseProxy
	fmt.Println("Hello from reverseProxy, I'm starting now!")
	fmt.Printf("I'm forwarding based on %s algorithm.\n", reverseProxy.Strategy)

	// Health checking 
	go HealthCheckFunc(reverseProxy, reverseProxy.HealthCheckFreq)

	//whenever a request is sent call reverseProxy.ServeHTTP
	http.Handle("/", reverseProxy)
	http.ListenAndServe(":" + strconv.Itoa(reverseProxy.Port), reverseProxy)

}
