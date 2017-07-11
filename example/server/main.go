package main

import (
	"github.com/wenzhenxi/tcp.io"
	"fmt"
)

func main() {
	tcp, err := tcpio.NewServer()

	if err != nil {
		panic(err)
	}

	tcp.On("connection", func(so tcpio.Socket) {
		fmt.Println("test")
	})

	tcp.Run("127.0.0.1", 9696)
}
