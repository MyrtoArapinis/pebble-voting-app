package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	prog := os.Args[0]
	cmd := os.Args[1]
	args := os.Args[2:]
	switch cmd {
	case "setup":
		if len(args) != 2 {
			fmt.Printf("Usage: %s setup DEPTH FILE\n", prog)
			return
		}
		depth, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		params, err := SetupCircuit(depth)
		if err != nil {
			fmt.Printf("Error creating params: %s\n", err.Error())
			return
		}
		bytes, err := params.ToBytes()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		filename := args[1]
		file, err := os.Create(filename)
		defer file.Close()
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
		_, err = file.Write(bytes)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return
		}
	}
}
