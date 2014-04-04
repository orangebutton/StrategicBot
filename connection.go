package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
)

type SBRequest map[string]interface{}

// Needs to look up the lobby URL, get the login token and then connect.
func Connect() {
	con, err := net.Dial("tcp", "107.21.58.31:8081")
	deny(err)

	log.Println("Connection established:", con)

	marshaledReq, err := json.Marshal(SBRequest{"msg": "LobbyLookup"})
	deny(err)

	log.Println("Marshaled Request:", marshaledReq)

	bytesWritten, err := con.Write(marshaledReq)
	deny(err)

	log.Println("Bytes written: ", bytesWritten)

	readBuffer := make([]byte, 1024)
	bytesRead, err := con.Read(readBuffer)
	deny(err)

	log.Println("Read:", readBuffer)
	log.Println("Bytes read:", bytesRead)

	var replyBuffer bytes.Buffer
	replyBuffer.Write(readBuffer[:bytesRead])

	log.Println("replyBuffer:", replyBuffer)
}