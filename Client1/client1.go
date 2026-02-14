package main

import (
	"io"
	"net/http"
	"log"
)

func main() {
	//Client sends requests to the Reverse Proxy
	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	//The client reads the response body returned by the reverse proxy 
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
}
