package main

import (
	"fmt"
	"net/http"
)



func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from backend server1", )
    })

    fmt.Println("Server running on :8081")
    http.ListenAndServe(":8081", nil)
}

