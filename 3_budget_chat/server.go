package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
)

var clients = make(map[string]connection)
var leaving = make(chan message)
var messages = make(chan message)

var validName = regexp.MustCompile("^[a-zA-Z0-9]+$")
var containsChar = regexp.MustCompile("[a-zA-Z]")

type message struct {
	text    string
	address string
}

type connection struct {
	conn net.Conn
	name string
}

func main() {
	addr := "0.0.0.0:9999"
	server, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer server.Close()

	log.Println("Server is running on:", addr)

	go broadcaster()
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Println("Failed to accept conn.", err)
			continue
		}

		go handler(conn)
	}
}

func handler(conn net.Conn) {
    defer func() {
        // Use SetLinger to force close the connection
        err := conn.(*net.TCPConn).SetLinger(0)
        if err != nil {
                log.Printf("Error when setting linger: %s", err)
        }
        log.Printf("Closing connection to " + conn.RemoteAddr().String())
        if err := conn.Close(); err != nil {
            log.Printf("Closing client %v connection returned error: %v\n", conn.RemoteAddr().String(), err)
        }
        log.Printf("Closed")
    }()


    log.Println("Accepted connection from: " + conn.RemoteAddr().String())
	fmt.Fprintln(conn, "Welcome to DamoChat! What shall I call you?")
	input := bufio.NewScanner(conn)

    input.Scan()
	name := input.Text()

    log.Println(conn.RemoteAddr().String() + " requested name: " + name)
	// Check name is legal!
	// must contain at least 1 char, and be only uppercase, lowercase, and digits
	if !isValidName(name) {
        log.Println(conn.RemoteAddr().String() + " requested invalid name: " + name)
		return
	}

    log.Println(conn.RemoteAddr().String() + " sending user list")
	fmt.Fprintln(conn, "* The room contains: "+allUsersNames())

	// Record new client
	// TODO: Needs to happen after getting name
	clients[conn.RemoteAddr().String()] = connection{
		conn: conn,
		name: name,
	}

	messages <- newMessage("* "+name+" has joined.", conn)

	for input.Scan() {
		messages <- newMessage("["+name+"] "+input.Text(), conn)
	}

	//Delete client form map
	delete(clients, conn.RemoteAddr().String())

	leaving <- newMessage("* "+name+" has left.", conn)
}
func isValidName(name string) bool {
	return containsChar.MatchString(name) && validName.MatchString(name)
}
func allUsersNames() string {
	var sb strings.Builder
	for _, conn := range clients {
		sb.WriteString(conn.name + ", ")
	}
	return sb.String()
}
func newMessage(msg string, conn net.Conn) message {
	addr := conn.RemoteAddr().String()
	return message{
		text:    msg,
		address: addr,
	}
}
func broadcaster() {
	for {
		select {
		case msg := <-messages:
			for _, conn := range clients {
				if msg.address == conn.conn.RemoteAddr().String() {
					continue
				}
				fmt.Fprintln(conn.conn, msg.text)
			}

		case msg := <-leaving:
			for _, conn := range clients {
				fmt.Fprintln(conn.conn, msg.text)
			}
		}
	}
}
