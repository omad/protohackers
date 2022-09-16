package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

type Request struct {
	Method *string  `json:"method"`
	Number *float64 `json:"number"`
}

type Response struct {
	Method string `json:"method"`
	Prime  bool   `json:"prime"`
}

func main() {
	addr := "0.0.0.0:9999"
	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer server.Close()

	log.Println("Server is running on:", addr)

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Failed to accept conn.", err)
			continue
		}
		log.Println("New connection from", conn.RemoteAddr())

		go handler(conn)
	}
}

func handler(conn net.Conn) error {
	defer conn.Close()

	c := bufio.NewReader(conn)
	table := make(map[int32]int32)
	for {
		msgtype, err := c.ReadByte()
		if err != nil {
			return err
		}

		var x, y int32
		err = binary.Read(c, binary.BigEndian, &x)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}
		err = binary.Read(c, binary.BigEndian, &y)
		if err != nil {
			fmt.Println("binary.Read failed:", err)
		}

		// fmt.Println(msgtype, x, y)
		if msgtype == 'I' {
			table[x] = y
		}
		if msgtype == 'Q' {
			sum := 0
			length := 0

			for k, v := range table {
				if x <= k && k <= y {
					length++
					sum += int(v)
				}
			}

			var mean int32
			if length != 0 {
				mean = int32(sum / length)
			}
			fmt.Println("Computed mean", mean, "length:", length, "tablesize", len(table))

			// send reply
			binary.Write(conn, binary.BigEndian, mean)
		}

	}

}
