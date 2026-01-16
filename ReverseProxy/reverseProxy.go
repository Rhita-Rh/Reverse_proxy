package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)


func main() {
	//Reverse Proxy will forward to one server in port :8081
    targetURL, err := url.Parse("http://localhost:8081")
	if err!=nil{
		log.Fatal()
	}

	//Our reverse Proxy that may be customized after  
	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
	
	//Strating reverseProxy
	fmt.Println("Hello from reverseProxy, I'm starting now!")
	fmt.Printf("I'm forwarding to %s", targetURL)

	http.Handle("/", reverseProxy)
	http.ListenAndServe(":8080", nil)
}
