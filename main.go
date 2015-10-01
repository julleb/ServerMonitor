package main


//  xmllint --valid --noout file.xml

import (
	"bytes"
	"encoding/xml"
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
	"time"
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
}

//handle the Server monitor page
func serverMonitorHandler(res http.ResponseWriter, req *http.Request) {
	//getting the ip from the url
	urlArray := strings.Split(req.URL.Path, "/")
	ip := urlArray[len(urlArray)-1]
	xmlString := getInformationFromDB(ip) //returns the old data as xml
	//here we can get the ip and query the db

    //fmt.Println(dtdValid(xmlString))
    if(!dtdValid(xmlString)) {
        //xml is not valid - but dont what we should do then ??
        fmt.Println("ERROR: xml is not valid...")
    }

    xslFile := determineStylesheet(req.UserAgent())

	htmlCode := processXSLTstdin(xslFile, xmlString) //"xslt-fake.xsl"
	io.WriteString(res, string(htmlCode))

}

//determine which stylesheet to use, depending on the users UserAgent
//in this way, we can use different stylesheet depending on the users platform
func determineStylesheet(userAgent string) (string) {
    userAgent = strings.ToLower(userAgent)
    if(strings.Contains(userAgent, "android")) {
        return "information-android-html.xsl"
    }
    return "information-html.xsl"
}


//check if the xml is valid to the dtd
func dtdValid(xml string) (bool) {
    b := bytes.NewBufferString(xml)
    cmd := exec.Cmd{
        Args: []string{"xmllint", "--valid", "--noout", "-"}, //"-"
        Env: os.Environ(),
        Path: "/usr/bin/xmllint",
        Stdin: b,
    }
    _,err:= cmd.Output()
    if err != nil {
        //its invalid
        return false    
    }
    return true
    //xmllint --valid --noout fake.xml
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

func processXSLTstdin(xslFile string, xmlString string) []byte {
	b := bytes.NewBufferString(xmlString)
	cmd := exec.Cmd{
		Args:  []string{"xsltproc", xslFile, "-"},
		Env:   os.Environ(),
		Path:  "/usr/bin/xsltproc",
		Stdin: b, // io.Reader
	}
	output, _ := cmd.Output()
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
        fmt.Println(messageFromInfoServer)
		if error != nil {
            //couldnt connect to the ip
            message = createMessage("-1") //error code -1, if clients get -1 he should know that the server doesnt exist
            err = conn.WriteMessage(messageType, message)
		    if err != nil {
			    fmt.Println(err)
			    return
		    }
			return
		}
        //we got a connect to the InfoServer
        
		//we have to insert the ip if it doesnt exists. Since we are doing a while loop we
		//dont want to query the db all the time, instead we use an variable for the checking
		if ipExist == false {
			if !ipExists(ip) { //the ip doesnt exist in db
				insertIP(ip)
			}
			ipExist = true
		}
        //insert the data into the db
		insertXMLtoDB(messageFromInfoServer, ip)

		//send the message to the firefox client
		message = createMessage(messageFromInfoServer)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

//Gets the xml data from the InfoServer via socket and return it as a string
func getDataFromInfoServer(ip string) (string, error) {
	ipAndPort := ip + ":9090"
	conn, err := net.Dial("tcp", ipAndPort)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	reply := make([]byte, 1024)
	conn.Read(reply)
	message := convertByteArrayToString(reply)
	return message, err
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
	if err != nil {
		fmt.Println(err)
	}

	var values []interface{}
	values = getDataFromXML(info.CPU.ServerData, values)
	values = getDataFromXML(info.Memory.ServerData, values)
	values = append(values, info.Date.Value)
	//lets insert the data  into the database
	insertInformation(ip, values)

}

//adds the average, min and max of cpu_temp to the struct holder from the db
func getTresholdsForCPU(holder *informations, ip string) {
	var values []interface{}
	values = append(values, ip)
	rows := db.Query("SELECT max(cpu_temp), min(cpu_temp), avg(cpu_temp) FROM server NATURAL JOIN has NATURAL JOIN information WHERE server.ip=$1", values)
	var max, min int
	var avg float32
	for rows.Next() {
		rows.Scan(&max, &min, &avg)
		u := Unit{Value: "C&degree;"}
		holder.Funfacts = funfacts{Attr: "temp", Min: min, Max: max, Avg: avg, Unit: u}

	}
}

//query the db to get all data from a certain ip addrs
func getInformationFromDB(ip string) string {
	dateLength := 19 //a date is always 19 letters long
	var holder informations
	var values []interface{}
	values = append(values, ip)
	rows := db.Query("SELECT * FROM server NATURAL JOIN has NATURAL JOIN information WHERE server.ip=$1", values)
	for rows.Next() {
		var info_id, cpu_temp, cpu_load, memory_usage, memory_total int
		var ip string
		var sqlDate time.Time
		rows.Scan(&info_id, &ip, &cpu_temp, &cpu_load, &memory_usage, &memory_total, &sqlDate)
		//since our sqlDate has some trash in the end, we need to remove it
		date := sqlDate.String()[:dateLength]
		// Javascript wants a T between date and time: 2012-01-01T12:00:00
		date = strings.Replace(date, " ", "T", 1)
		holder = dataToXML(holder, info_id, ip, cpu_temp, cpu_load, memory_usage, memory_total, date)
	}

	getTresholdsForCPU(&holder, ip)

	db.DeferRows(rows)
	//convert the holder to XML
	output, err := xml.MarshalIndent(holder, "", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
    //getXMLHeader() + convertByteArrayToString(output)
	return getXMLHeader() + convertByteArrayToString(output)
}


//gets the data from the xml and puts it in the values array
func getDataFromXML(serverdata []serverData, values []interface{}) []interface{} {
	for i := 0; i < len(serverdata); i++ {
		values = append(values, serverdata[i].Value)
	}
	return values
}

//insert data into the information table
//and create an relation between the information and the ip in the db
func insertInformation(ip string, values []interface{}) {
	//would be nice to do a transaction here, for the coolness
	rows := db.Query("INSERT INTO information(cpu_temp,cpu_load,memory_usage,memory_total, date) VALUES($1,$2,$3,$4,$5) RETURNING info_id", values)
	var info_id int
	for rows.Next() {
		rows.Scan(&info_id)
	}
	var hasValues []interface{}
	hasValues = append(hasValues, ip)
	hasValues = append(hasValues, info_id)
	rows = db.Query("INSERT INTO has(ip, info_id) VALUES($1,$2)", hasValues)
	db.DeferRows(rows)
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
		db.DeferRows(rows)
		return true
	}
	db.DeferRows(rows)
	return false
}
