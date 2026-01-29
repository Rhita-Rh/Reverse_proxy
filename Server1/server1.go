package main

import (
	"fmt"
	"net/http"
	"time"
)



func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(3 * time.Second)
        fmt.Fprintf(w, "Hello from backend server1\n")
    })

    fmt.Println("Server running on :8081")
    http.ListenAndServe(":8081", nil)
}

