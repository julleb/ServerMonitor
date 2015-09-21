package main

import (
    "fmt"
    "net/http"
    "io"
)



func main() {
    http.HandleFunc("/", hello)
    fmt.Println("This server is going up on port 8080")
    http.ListenAndServe(":8080", nil)
    
}

func hello(res http.ResponseWriter, req *http.Request) {
    io.WriteString(res, "Hello my friend")
    
    
}


