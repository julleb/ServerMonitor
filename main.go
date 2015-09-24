package main

import (
	"fmt"   
	"html/template"
	"io"
	"net/http"
	"os"
	"os/exec"
    "net"
    db "github.com/julleb/DBFuncs"
    "github.com/gorilla/websocket"
    "github.com/gorilla/mux"
    "strings"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024, //might need to increase this?
    WriteBufferSize: 1024,
}

type temp struct {
	X string
	Y string
}

//compile all templates and cache them
var templates = template.Must(template.ParseGlob("views/*"))

func main() {
    
    db.OpenDBConnection()
    
    r := mux.NewRouter()
    r.HandleFunc("/public/", visualHandler);
	r.HandleFunc("/", index)
    
    r.HandleFunc("/requestdata/{ip}", requestDataHandler)
    r.HandleFunc("/newip", formHandler)
	r.HandleFunc("/{ip}", serverMonitorHandler)
    http.Handle("/", r)    
    
    /*
	http.HandleFunc("/public/", visualHandler)
	http.HandleFunc("/", index)
    
    http.HandleFunc("/requestdata", requestDataHandler)
	http.HandleFunc("/newip", formHandler)
    */

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
	ip := req.PostFormValue("ip")
	fmt.Println(ip)
	//do some cool stuffs
	//check if ip is in db
    if ipExists(ip) {
        //ip exists in the db        
        fmt.Println("IP DOES EXIST")
        //need to query to get last week info        
    }
    //redirect the user to the ip url
    http.Redirect(res, req, "/"+ip, 301)
	//templates.ExecuteTemplate(res, "index", s) //render a page
}

//handle the Server monitor page
func serverMonitorHandler(res http.ResponseWriter, req *http.Request) {
    //getting the ip from the url    
    urlArray := strings.Split(req.URL.Path, "/")
    ip := urlArray[len(urlArray)-1]
    
    if ipExists(ip) {
       var values []interface{}
       values = append(values, ip)
       rows := db.Query("SELECT * FROM server NATURAL JOIN has NATURAL JOIN information WHERE server.ip=has.ip", values);
       for rows.Next() {
            

       }      
    }
    //here we can get the ip and query the db
    htmlCode := processXSLT("xslt-fake.xsl", "fake.xml")
    io.WriteString(res, string(htmlCode))

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

func requestDataHandler(res http.ResponseWriter, req *http.Request) {    
    
    //getting the ip from the url    
    urlArray := strings.Split(req.URL.Path, "/")
    ip := urlArray[len(urlArray)-1]
      
    conn, err := upgrader.Upgrade(res, req, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    //while loop
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            //the user disconnected
            fmt.Println(err)
            return
        }
        
        //the message from the server
        messageFromInfoServer := getDataFromInfoServer(ip)

        
        fmt.Println("i got this " ,messageFromInfoServer)

        //send the message to the firefox client
        message = createMessage("i can be your hero baby")
        err = conn.WriteMessage(messageType, message);
        if  err != nil {
            return
        }
    }
}

func getDataFromInfoServer(ip string) (string) {
    ipAndPort := ip + ":9090"
    conn, err := net.Dial("tcp", ipAndPort)
    if(err != nil) {
        fmt.Println(err)
        return "";  
    }
    reply := make([]byte, 1024)
    conn.Read(reply)
    message := convertByteArrayToString(reply)
    //_,_ = bufio.NewReader(conn).ReadString('\n');
    
    return message
}

func createMessage(message string) ([]byte) {
    return []byte(message)
}

func convertByteArrayToString(arr []byte) (string) {
    return string(arr[:])
}

func insertIP(ip string) {
    var values []interface{}
    values = append(values, ip)
    row := db.Query("INSERT INTO server(ip) values($1)", values )
    //db.DeferRows(row)    
    fmt.Println(row.Columns())
    
    for row.Next() {
        var col string
        row.Scan(&col)
        fmt.Println("Weee " + col)
    }
    db.DeferRows(row)
}




func insertInformation(values []interface{}) {
   _ = db.Query("INSERT INTO information(id,cpu_temp,memory_usage,memory_total,total_memory) VALUES($1, $2,$3,$4,$5)", values)     
}

func ipExists(ip string) (bool) {
    var values []interface{}
    values = append(values, ip)
    rows := db.Query("SELECT * FROM SERVER WHERE IP=$1", values)
    for rows.Next() {
        return true
    }
    return false
}



