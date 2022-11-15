package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
)

var database = make(map[string]string)

func main() {
	addr := "0.0.0.0:9999"
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer pc.Close()

	log.Println("Server is running on:", addr)

	for {
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		serve(pc, addr, buf[:n])
	}
}

func serve(pc net.PacketConn, addr net.Addr, buf []byte) {
	log.Printf("Received packet '%s'", buf)
	before, after, found := bytes.Cut(buf, []byte("="))
	if found {
		// This is an Insert message
		log.Println("Received insert message")
		if string(before) == "version" {
			log.Println("Cannot insert to version")
			return
		}
		database[string(before)] = string(after)
	} else {
		log.Printf("Received query for %s", buf)
		if string(buf) == "version" {
			pc.WriteTo([]byte("version=Damo's dodgy KV store"), addr)
		} else {
			resp := []byte(fmt.Sprintf("%s=%s", buf, database[string(buf)]))
			log.Printf("Responding with '%s'", resp)
			pc.WriteTo(resp, addr)
		}

	}

}
