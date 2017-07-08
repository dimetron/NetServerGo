package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"net"
	"os"
	"sync/atomic"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

type GpsDeviceEvent struct {
	UID uint32
	TXT string
}

//global counter
var requestIdCounter uint32 = 0

func main() {
	//mongo session
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("carmonit").C("gpsdata")
	c.RemoveAll(nil)

	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, session)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, session *mgo.Session) {
	defer conn.Close()

	addr := conn.RemoteAddr()

	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	//insert mongo
	txt := string(buf[:reqLen])
	insertMongo(session, txt)

	if err != nil {
		log.Fatal(err)
		conn.Write([]byte("ERR"))
		fmt.Println("Failed to process: [", reqLen, "] from", addr.String())
	} else {
		// Send a response back to person contacting us.
		conn.Write([]byte("OK"))
		fmt.Println("Received bytes: [", reqLen, "] from", addr.String())
	}
}

func insertMongo(session *mgo.Session, txt string) (err error) {
	UID := atomic.AddUint32(&requestIdCounter, 1)

	c := session.DB("carmonit").C("gpsdata")
	err = c.Insert(&GpsDeviceEvent{UID, txt})
	if err != nil {
		log.Fatal(err)
	}
	return err
}
