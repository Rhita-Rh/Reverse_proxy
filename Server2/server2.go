package main

import (
	"fmt"
	"net/http"
	"time"
)



func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(3 * time.Second)
        fmt.Fprintf(w, "Hello from backend server2\n")
    })

    fmt.Println("Server running on :8082")
    http.ListenAndServe(":8082", nil)
}

