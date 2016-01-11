package main

import (
	"encoding/xml"
)

type Unit struct {
	XMLName xml.Name `xml:"Unit"`
	Value   string   `xml:",innerxml"`
}

type serverData struct {
	Description string
	Value       int    `xml:"value"`
	Unit        Unit   `xml:"Unit"`
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

type Date struct {
	XMLName xml.Name `xml:"Date"`
	Value   string   `xml:"Date"`
}

type information struct {
	XMLName xml.Name `xml:"information"`
	Date
	CPU
	Memory
}
type funfacts struct {
	Attr string  `xml:",attr"`
	Min  int     `xml:"Min"`
	Max  int     `xml:"Max"`
	Avg  float32 `xml:"Avg"`
	Unit Unit
}

type informations struct {
	XMLName  xml.Name `xml:"informations"`
	Infos    []information
	Funfacts funfacts
}

const header = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE informations SYSTEM "information.dtd">`

/*`<?xml version="1.0" encoding="UTF-8"?>


<!DOCTYPE informations [
<!ELEMENT informations (information*, Min, Max) >
<!ELEMENT information (Date, CPU*, Memory*) >
<!ELEMENT Date (#PCDATA) >
<!ELEMENT CPU (ServerData*)>
<!ELEMENT Memory (ServerData*)>
<!ELEMENT ServerData (Description, value)>
<!ATTLIST ServerData unit CDATA #IMPLIED>
<!ELEMENT Description (#PCDATA)>
<!ELEMENT value (#PCDATA)>
<!ELEMENT Max (#PCDATA)>
<!ELEMENT Min (#PCDATA)>

<!ENTITY degree "&#176;">

]>
`*/

//returns the xml header
func getXMLHeader() string {
	return header
}

//creates and xml struct with the given parameters and returns the struct.
//we can use this struct to generate and xml string
func dataToXML(holder informations, info_id int, ip string, cpu_temp int, cpu_load int, memory_usage int, memory_total int, sqlDate string) informations {

	t := &serverData{Description: "Temperature", Value: cpu_temp, Unit: Unit{Value: "C&degree;"}}

	t1 := &serverData{Description: "Load", Value: cpu_load, Unit: Unit{Value: "%"}}
	a := []serverData{*t, *t1}

	f := &serverData{Description: "Total", Value: memory_total, Unit: Unit{Value: "MB"}}
	f1 := &serverData{Description: "Used", Value: memory_usage, Unit: Unit{Value: "MB"}}
	b := []serverData{*f, *f1}

	cpu := &CPU{ServerData: a}
	mem := &Memory{ServerData: b}
	date := &Date{Value: sqlDate}
	i := information{Date: *date, CPU: *cpu, Memory: *mem}
	holder.Infos = append(holder.Infos, i)
	return holder
}
