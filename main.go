package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
)

type temp struct {
	X string
	Y string
}

//compile all templates and cache them
var templates = template.Must(template.ParseGlob("views/*"))

func main() {
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

	//we just test the db here
	dbConnection()
    
    //we test insertion
    insertInformation()

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
}

//insert function to db
//TODO we have to make variable db global, so we dont make a connection do db all the time
func insertInformation() {
	db, err := sql.Open("postgres", "user=postgres password=lol dbname=servermonitor")
    ip := "50"
	stmt, err := db.Prepare("insert into server(ip) values($1)")
	if err != nil {
		
	}
    //defer it, very important, to avoid runtime panic
	defer stmt.Close()
	rows, err := stmt.Query(ip)
	if err != nil {
		
	}
    //defer it, very important, to avoid runtime panic
	defer rows.Close()
	
}
