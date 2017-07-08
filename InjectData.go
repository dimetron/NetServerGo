package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

//865733021674619,cSQ88qShwC3,08/09/15,15:47:57+0,80915,15475800,50.097645,14.436432,51.78,185.10,244.27,74,16,14.17,1

type GpsLog struct {
	IMEI string
	CODE string
	DATE string
	TIME string
	DTST string
	TMST string
	LAT  string
	LON  string
	SPD  string
	DIR  string
}

func main() {
	concurrency := 5
	sem := make(chan bool, concurrency)

	file, err := os.Open("testdata/inputs.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() { // internally, it advances token based on separator
		i++
		line := scanner.Text()
		parts := strings.Split(line, ",")

		gs := &GpsLog{
			IMEI: parts[0],
			CODE: parts[1],
			DATE: parts[2],
			TIME: parts[3],
		}

		sem <- true
		go func(i int, gs *GpsLog, ln string) {
			defer func() { <-sem }()
			sendTCP(i, gs, ln)
		}(i, gs, line)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func sendTCP(i int, gs *GpsLog, txt string) {

	servAddr := "localhost:3333"
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(txt))
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	reply := make([]byte, 1024)

	_, err = conn.Read(reply)
	if err != nil {
		println("Write to server failed:", err.Error())
		os.Exit(1)
	}

	fmt.Print("\n* ", i, gs, " -> ", string(reply))
}
