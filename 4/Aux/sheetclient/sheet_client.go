package sheetclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

const NETWORK_PORT = ":8000"

// Run the execution flow for the sheet client
func Run(server, name string) {
	// Establish TCP connection to server
	r, w, err := connectTCP(server)
	if err != nil {
		panic(err)
	}

	// Make initial handshake request
	identifier := initialRequest(w, r, name)
	fmt.Println("Client assigned ID: ", identifier)

	// Read from stdin
	decoder := json.NewDecoder(os.Stdin)

	for {
		jsonData := readJSON(decoder)
		response := sendJSON(r, w, jsonData)
		fmt.Println("Response:", response)
	}
}

// Connect to a TCP Server and return a reader
func connectTCP(server string) (io.Reader, io.Writer, error) {
	conn, err := net.Dial("tcp", server+NETWORK_PORT)
	if err != nil {
		return nil, nil, errors.New("Dialing " + server + " failed: " + err.Error())
	}
	return conn, conn, nil
}

// Send a UUID, receive a unique internal identifier
func initialRequest(w io.Writer, r io.Reader, name string) string {

	byteName, _ := json.Marshal(name)
	fmt.Println("Sending name: ", name)
	w.Write(byteName)

	fmt.Println("Waiting for ID")
	buf := make([]byte, 1024)
	if n, err := r.Read(buf); err != nil {
		panic(err)
	} else {
		return string(buf[:n])
	}
}

// Connect to a TCP Server and return a reader
func readJSON(decoder *json.Decoder) []interface{} {
	returns := make([]interface{}, 0)
	var jsonVal interface{}
	for decoder.More() {
		decoder.Decode(&jsonVal)
		returns = append(returns, jsonVal)
		if isAtRequest(jsonVal) {
			break
		}
	}
	return returns
}

// Decipher if an interface is an "at" request
func isAtRequest(req interface{}) bool {
	arr, isArr := req.([]interface{})
	if !isArr {
		return false
	}

	if len(arr) < 1 {
		return false
	}

	if str, isStr := arr[0].(string); !isStr {
		return false
	} else {
		return str == "at"
	}
}

// Send data over a read/writer and return the string received
func sendJSON(r io.Reader, w io.Writer, data []interface{}) string {

	fmt.Println("Data:", data)

	marshalled, _ := json.Marshal(data)
	fmt.Println("Sending:", string(marshalled))

	w.Write(marshalled)
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	return string(buf[:n])
}
