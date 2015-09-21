package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type temp struct {
	X int
	Y int
}

//compile all templates and cache them
var templates = template.Must(template.ParseGlob("views/*"))

func main() {
	http.HandleFunc("/public/", visualHandler)
	http.HandleFunc("/", hello)

	fmt.Println("This server is going up on port 8080")
	http.ListenAndServe(":8080", nil)

}

func hello(res http.ResponseWriter, req *http.Request) {
    fmt.Println("im in hello")
	s := temp{1, 2}
	templates.ExecuteTemplate(res, "index", s)

}

func visualHandler(res http.ResponseWriter, req *http.Request) {
    fmt.Println("some1 asked for ", req.URL.Path[1:])	
    http.ServeFile(res, req, req.URL.Path[1:])
    
}
