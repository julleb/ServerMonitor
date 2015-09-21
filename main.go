package main

import (
    "fmt"
    "net/http"
    "html/template"
)


type temp struct {
    X int
    Y int
}

//compile all templates and cache them
var templates = template.Must(template.ParseGlob("views/*"))

func main() {
    http.HandleFunc("/", hello)
    fmt.Println("This server is going up on port 8080")
    http.ListenAndServe(":8080", nil)
    
}

func hello(res http.ResponseWriter, req *http.Request) {
    s := temp{1,2}
    templates.ExecuteTemplate(res, "index", s)
    
}


