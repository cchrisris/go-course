package main

import (
	"encoding/base64"
	"fmt"

	"github.com/cchrisris/go-course/hw2/internal/client"
)

var ServerUrl = "http://127.0.0.1:8080"

func main() {
	newClient := client.NewClient(ServerUrl)

	version, err := newClient.GetVersion()
	if err != nil {
		panic(err)
	}
	fmt.Println(version)

	originalString := "MTS Go Homework #2"
	decoded, err := newClient.DecodeString(base64.StdEncoding.EncodeToString([]byte(originalString)))
	if err != nil {
		panic(err)
	}
	fmt.Println(decoded)

	_, responseCode, _, err := newClient.GetHardOp()
	if err != nil {
		fmt.Printf("%t, %d\n", false, 0)
		return
	}
	fmt.Printf("%t, %d\n", true, responseCode)
}
