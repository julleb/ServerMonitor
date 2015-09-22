package main

//https://systembash.com/a-simple-go-tcp-server-and-tcp-client/

import (
	"bufio"
	"fmt"
	"net"
    "strings"
)

func main() {
    
    requestString := "REQ123"

	fmt.Println("======= MoniGtor server =======")

	listener, _ := net.Listen("tcp", ":9090")
	conn, _ := listener.Accept()
    
    recievedMessage,_ := bufio.NewReader(conn).ReadString('\n')
    
    if(Compare(recievedMessage, requestString) == 0) {
        //they request the data
        sendMessage := "hejsan"
	    conn.Write([]byte(sendMessage))
    }else {
        //they didnt request any data, ignore them?
    }
    

}
