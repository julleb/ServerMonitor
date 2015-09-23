package main

import (
	"fmt"   
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
    db "github.com/julleb/DBFuncs"
)

type temp struct {
	X string
	Y string
}

//compile all templates and cache them
var templates = template.Must(template.ParseGlob("views/*"))

func main() {
    
    db.OpenDBConnection()
    var rows db.Rows
    t := db.Tuple{"ip", "291.293.22.2"}
    rows.Tuples = append(rows.Tuples, t)
    db.InsertIntoTable("server", rows)
	http.HandleFunc("/public/", visualHandler)
	http.HandleFunc("/", index)
    
	http.HandleFunc("/newip", formHandler)

	fmt.Println("This server is going up on port 8080")
	http.ListenAndServe(":8080", nil)

}

// index page
func index(res http.ResponseWriter, req *http.Request) {
	fmt.Println("im in hello")
	s := temp{"JULLE", "DDANI"}



	templates.ExecuteTemplate(res, "index", s)

}

//function for css and js
func visualHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, req.URL.Path[1:])

}

//function for handling the html form
func formHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("some1 asked for formhandler ")
	ip := req.PostFormValue("ip")
	fmt.Println(ip)
	//do some cool stuffs
	//check if ip is in db
	//some db logic stuffs
	//s := temp{ip, ip}
	htmlCode := processXSLT("xslt-fake.xsl", "fake.xml")
	io.WriteString(res, string(htmlCode))

	//templates.ExecuteTemplate(res, "index", s) //render a page

}

func processXSLT(xslFile string, xmlFile string) []byte {
	cmd := exec.Cmd{
		Args: []string{"xsltproc", xslFile, xmlFile},
		Env:  os.Environ(),
		Path: "/usr/bin/xsltproc",
	}
	output, _ := cmd.Output()
	fmt.Printf("yooo %s\n", output)
	return output
}

/*
func dbConnection() {
	db, err := sql.Open("postgres", "user=postgres password=lol dbname=servermonitor")
	if err != nil {
		fmt.Println("No connectioN!")
	}
	fmt.Println("we have connection")
	//...where ip=$1", ip)
	rows, err := db.Query("SELECT * from information")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close() // defer it, very important, to avoid runtime panic

	for rows.Next() { // looping through all rows
		var id, cpu, us, to, too int
		rows.Scan(&id, &cpu, &us, &to, &too) // getting  all the cols value
		fmt.Println(id)
		//fmt.Println(rows.Columns())
	}
}*/



