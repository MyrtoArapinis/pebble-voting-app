package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/giry-dev/pebble-voting-app/pebble-core/server"
	"golang.org/x/term"
)

var flagPassHash = flag.String("passhash", "", "server password hash (sha256)")

func main() {
	var passHash []byte
	var err error
	flag.Parse()
	if *flagPassHash != "" {
		passHash, err = hex.DecodeString(*flagPassHash)
		if err != nil {
			fmt.Println("Error decoding password hash: ", err)
			return
		}
	}
	mode := flag.Arg(0)
	switch mode {
	case "hash":
		fmt.Print("Input: ")
		input, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return
		}
		fmt.Print("\nConfirm: ")
		confirm, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return
		}
		fmt.Println()
		if !bytes.Equal(input, confirm) {
			fmt.Println("Inputs don't match")
			return
		}
		output := sha256.Sum256(input)
		fmt.Println(hex.EncodeToString(output[:]))
	case "mock":
		endpoint := flag.Arg(1)
		handler := server.NewMockServer(endpoint, passHash)
		fmt.Println("Starting mock server...")
		err = http.ListenAndServe(endpoint, handler)
		if err != nil {
			fmt.Println(err)
		}
	}
}
