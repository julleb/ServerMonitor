package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	db "github.com/julleb/DBFuncs"
	"html/template"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
    "encoding/xml"
)



type serverData struct {
	Description string
	Value       int    `xml:"value"`
	Unit        string `xml:"unit,attr"`
	Comment     string `xml:",comment"`
}

type CPU struct {
	XMLName    xml.Name     `xml:"CPU"`
	ServerData []serverData `xml:"CPU>ServerData"`
}

type Memory struct {
	XMLName    xml.Name     `xml:"Memory"`
	ServerData []serverData `xml:"Memory>ServerData"`
}

type information struct {
	XMLName xml.Name `xml:"information"`
	CPU
	Memory
}


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
	//r.HandleFunc("/public/", visualHandler)
	r.HandleFunc("/", index)

	r.HandleFunc("/requestdata/{ip}", requestDataHandler)
	r.HandleFunc("/newip", formHandler)
	r.HandleFunc("/{ip}", serverMonitorHandler)

	s := http.StripPrefix("/public/", http.FileServer(http.Dir("./public/")))
	r.PathPrefix("/public/").Handler(s)
	http.Handle("/", r)
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
	fmt.Println("serving:", req.URL.Path[1:])
	http.ServeFile(res, req, req.URL.Path[1:])

}

//function for handling the html form
func formHandler(res http.ResponseWriter, req *http.Request) {
	ip := req.PostFormValue("ip")
	//redirect the user to the ip url
	http.Redirect(res, req, "/"+ip, 301)
	//templates.ExecuteTemplate(res, "index", s) //render a page
}

//handle the Server monitor page
func serverMonitorHandler(res http.ResponseWriter, req *http.Request) {
	//getting the ip from the url
	urlArray := strings.Split(req.URL.Path, "/")
	ip := urlArray[len(urlArray)-1]
    getInformationFromDB(ip)
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
	//fmt.Printf("yooo %s\n", output)
	return output
}

func requestDataHandler(res http.ResponseWriter, req *http.Request) {    
    
    ipExist := false //used for not querying the db all the time
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
        messageFromInfoServer, error := getDataFromInfoServer(ip)
        if(error != nil) {
            //couldnt connect to the ip
            return
        }
        
        //we have to insert the ip if it doesnt exists. Since we are doing a while loop we
        //dont want to query the db all the time, instead we use an variable for the checking
        
         if ipExist == false {
            if !ipExists(ip) { //the ip doesnt exist in db
                insertIP(ip) 
            }
            ipExist = true  
        }

        insertXMLtoDB(messageFromInfoServer, ip)
        
        //we got a connect to the InfoServer
        //TODO add the data to the server!
               


        //send the message to the firefox client
        message = createMessage(messageFromInfoServer)
        err = conn.WriteMessage(messageType, message);
        if  err != nil {
            fmt.Println(err)
            return
        }
    }
}


//Gets the xml data from the InfoServer via socket and return it as a string
func getDataFromInfoServer(ip string) (string, error) {
    ipAndPort := ip + ":9090"
    conn, err := net.Dial("tcp", ipAndPort)
    if(err != nil) {
        fmt.Println(err)
        return "", err;  
    }
    reply := make([]byte, 1024)
    conn.Read(reply)
    message := convertByteArrayToString(reply)
    //_,_ = bufio.NewReader(conn).ReadString('\n');
    
    return message,err
}

//used for converting the message from string to byte
func createMessage(message string) []byte {
	return []byte(message)
}

//used for converting message from byte to string
func convertByteArrayToString(arr []byte) string {
	return string(arr[:])
}


func insertXMLtoDB(xmldata string, ip string) {
    info := information{} //CPU: none, Memory: none
    err := xml.Unmarshal([]byte(xmldata), &info)
    if(err != nil) {
        fmt.Println(err)
    }
    
    var values []interface{}
    values = getDataFromXML(info.CPU.ServerData, values)
    values = getDataFromXML(info.Memory.ServerData, values)
    //lets insert the data  into the database
    insertInformation(ip ,values)
    
    
}

//query the db to get all data from a certain ip addrs
func getInformationFromDB(ip string)  {
      var values []interface{}
      values = append(values, ip)
      rows := db.Query("SELECT * FROM server NATURAL JOIN has NATURAL JOIN information WHERE server.ip=$1", values);
      fmt.Println("-------")
      fmt.Println(rows.Columns())
      for rows.Next() {
            var info_id,cpu_temp, cpu_load, memory_usage, memory_total int
            var date, ip string
            rows.Scan(&info_id, &ip, &cpu_temp, &cpu_load, &memory_usage, &memory_total, &date)  
            fmt.Println(info_id)
            fmt.Println(cpu_temp)
            fmt.Println(memory_usage)
                 
      }
      fmt.Println("-------")
}

//gets the data from the xml and puts it in the values array
func getDataFromXML(serverdata []serverData, values []interface{}) ([]interface{}){ 
    for i := 0; i < len(serverdata); i++ {
        values = append(values,  serverdata[i].Value)
    }    
    return values
} 


//insert data into the information table 
//and create an relation between the information and the ip in the db
func insertInformation(ip string ,values []interface{}) {
    //would be nice to do a transaction here, for the coolness	
    rows := db.Query("INSERT INTO information(cpu_temp,cpu_load,memory_usage,memory_total, date) VALUES($1,$2,$3,$4, now()) RETURNING info_id", values)
    var info_id int
    for rows.Next() {
        rows.Scan(&info_id)
    }
    var hasValues []interface{}
    hasValues = append(hasValues, ip)
    hasValues = append(hasValues, info_id)
    rows = db.Query("INSERT INTO has(ip, info_id) VALUES($1,$2)", hasValues)
}

//addes the ip into the database
func insertIP(ip string) {
	var values []interface{}
	values = append(values, ip)
	rows := db.Query("INSERT INTO server(ip) values($1)", values)
    db.DeferRows(rows)
    
}
//checks if the ip exists in the database
func ipExists(ip string) bool {
	var values []interface{}
	values = append(values, ip)
	rows := db.Query("SELECT * FROM SERVER WHERE IP=$1", values)
	for rows.Next() {
		return true
	}
	return false
}
