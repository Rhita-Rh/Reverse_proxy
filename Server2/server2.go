package main

import (
	"fmt"
	"net/http"
)



func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello from backend server2", )
    })

    fmt.Println("Server running on :8082")
    http.ListenAndServe(":8082", nil)
}

