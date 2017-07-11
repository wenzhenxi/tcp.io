package main

import (
	"net"
	"fmt"
	"os"
)

func main() {

	conn, err := net.Dial("tcp", ":9696")
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}

	test := make([]byte,1024)
	conn.Read(test)

	fmt.Println(test)
}
