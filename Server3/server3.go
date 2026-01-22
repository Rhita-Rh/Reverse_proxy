package main

import (
	"fmt"
	"net/http"
)



func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from backend server3\n")
    })

    fmt.Println("Server running on :8083")
    http.ListenAndServe(":8083", nil)
}

