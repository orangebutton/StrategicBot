package main

import (
	"bytes"
	"encoding/json"
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
	log.Println("===== Connection =====")

	url := GetLobbyURL()
	token := GetLoginToken(email, password)

	log.Println(url, token)
}

// Listen to an URL and send line by line into a channel
// Returns the connection and said channel
func ListenToURL(url string) (net.Conn, chan []byte) {
	// Connect to the specified URL
	con, err := net.Dial("tcp", url)
	deny(err)

	log.Println("Listening on new connection:", url)

	// Make the channel (it can send and recieve byte-slices)
	ch := make(chan []byte)

	go func() {
		var chBuffer bytes.Buffer
		readFromCon := make([]byte, 1024)

		for {
			// Read 1024 bytes
			bytesRead, err := con.Read(readFromCon)
			if err != nil {
				close(ch)
				log.Println("Connection error:", err)
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
	log.Println("Get lobby URL...")
	con, ch := ListenToURL("107.21.58.31:8081")
	defer con.Close()
	defer log.Println("Connection closed: 107.21.58.31:8081")

	SendRequest(con, Request{"msg": "LobbyLookup"})

	for reply := range ch {
		var v MLobbyLookup
		json.Unmarshal(reply, &v)
		if v.Msg == "LobbyLookup" {
			url := v.Ip + ":" + strconv.Itoa(v.Port)
			log.Println("Lobby URL is", url)
			return url
		}
	}

	return ""
}

func GetLoginToken(email, password string) LoginToken {
	log.Println("Get login token...")

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

	log.Println("Recieved login token")

	return token
}