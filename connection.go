package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
)

type Request map[string]interface{}
type LoginToken map[string]interface{}

func SendRequest(con net.Conn, req Request) bool {
	marshaledReq, err := json.Marshal(req)
	deny(err)

	_, err = con.Write(marshaledReq)

	return err == nil
}

func Connect(email, password string) {
	log.Println("LobbyURL is", GetLobbyURL())
	log.Println("LoginToken is ", GetLoginToken(email, password))
}

// Listen to an URL and send line by line into a channel
// Returns the connection and said channel
func ListenToURL(url string) (net.Conn, chan []byte) {
	// Connect to the specified URL
	con, err := net.Dial("tcp", url)
	deny(err)

	// Make the channel (it can send and recieve byte-slices)
	ch := make(chan []byte)

	go func() {
		var chBuffer bytes.Buffer
		readFromCon := make([]byte, 1024)

		for {
			// Read 1024 bytes
			bytesRead, err := con.Read(readFromCon)
			if err == io.EOF {
				close(ch)
				log.Printf("Reached end of file. Connection closed.")
				return
			}

			// Buffer them
			chBuffer.Write(readFromCon[:bytesRead])

			// Cut into lines
			lines := bytes.SplitAfter(chBuffer.Bytes(), []byte("\n"))

			// Send lines to the through the channel
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

func GetLobbyURL() string {
	con, ch := ListenToURL("107.21.58.31:8081")
	defer con.Close()

	SendRequest(con, Request{"msg": "LobbyLookup"})

	for reply := range ch {
		var v MLobbyLookup
		json.Unmarshal(reply, &v)
		if v.Msg == "LobbyLookup" {
			return v.Ip + ":" + strconv.Itoa(v.Port)
		}
	}

	return ""
}

func GetLoginToken(email, password string) LoginToken {
	req := Request{
		"agent": Request{
			"name": "Scrolls",
			"version": 1,
		},
		"username": email,
		"password": password,
	}

	marshaledReq, err := json.Marshal(req)
	deny(err)

	buf := bytes.NewBufferString(string(marshaledReq))

	resp, err := http.Post("https://authserver.mojang.com/authenticate", "application/json", buf)
	deny(err)
	defer resp.Body.Close()

	readBuf := make([]byte, 2000)

	bytesRead, err := resp.Body.Read(readBuf)
	deny(err)

	var token LoginToken
	err = json.Unmarshal(readBuf[:bytesRead], &token)
	deny(err)

	return token
}