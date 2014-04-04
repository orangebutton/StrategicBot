package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"strconv"
)

type SBRequest map[string]interface{}

// Needs to look up the lobby URL, get the login token and then connect.
func Connect() {
	con, ch := ListonToURL("107.21.58.31:8081")

	log.Println("Connection established:", con)

	marshaledReq, err := json.Marshal(SBRequest{"msg": "LobbyLookup"})
	deny(err)

	log.Println("Marshaled Request:", marshaledReq)

	bytesWritten, err := con.Write(marshaledReq)
	deny(err)

	log.Println("Bytes written: ", bytesWritten)

	for reply := range ch {
		var v MLobbyLookup
		json.Unmarshal(reply, &v)
		if v.Msg == "LobbyLookup" {
			lobbyURL := v.Ip + ":" + strconv.Itoa(v.Port)
			log.Println("v is:", v)
			log.Println("lobbyURL is:", lobbyURL)
		}
	}
}

// Listens to an URL and buffers the read bytes.
func ListonToURL(url string) (net.Conn, chan []byte) {
	// Connect to the specified URL
	con, err := net.Dial("tcp", url)
	deny(err)

	// Make a channel to pass the listened content to other functions
	ch := make(chan []byte)

	// Goroutine which listens and cuts the content into convenient portions
	go func() {
		var chBuffer bytes.Buffer
		readBuffer := make([]byte, 1024)

		for {
			bytesRead, err := con.Read(readBuffer)
			if err != nil {
				close(ch)
				log.Printf("ListenToURL connection error: %s", err)
				return
			}

			chBuffer.Write(readBuffer[:bytesRead])

			// Cut into lines
			lines := bytes.SplitAfter(chBuffer.Bytes(), []byte("\n"))

			for _, line := range lines[:len(lines)-1] {
				n := len(line)
				if n > 1 {
					lineCopy := make([]byte, n)
					copy(lineCopy, line)
					ch <- lineCopy
				}

				chBuffer.Next(n)
			}
		}
	}()

	return con, ch
}