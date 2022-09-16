package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"math/big"
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

		go handler(conn)
	}
}

func handler(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		req, err := parseRequest(scanner.Bytes())

		if err != nil {
			// invalid json
			conn.Write([]byte(`{"response": "` + err.Error() + `"}` + "\n"))
			conn.Close()
			return

		}

		result := IsPrime(*req.Number)

		resp := Response{"isPrime", result}
		data, _ := json.Marshal(resp)

		conn.Write(append(data, '\n'))
	}
	if scanner.Err() != nil {
		log.Println("Error: ", scanner.Err())
	}
}
func parseRequest(request []byte) (*Request, error) {
	req := Request{}
	err := json.Unmarshal(request, &req)

	log.Println("Received Line:", string(request))

	if err != nil {
		log.Println("Invalid JSON:", string(request))
		return nil, errors.New("invalid json")
	}

	if req.Method == nil || req.Number == nil {
		log.Println("Received Incomplete Request")
		return nil, errors.New("incomplete request")
	}

	log.Println("Received Request: ", *req.Method, *req.Number)
	if *req.Method != "isPrime" {
		log.Println("Received Invalid Method")
		return nil, errors.New("invalid method")
	}
	return &req, nil
}

func IsPrime(num_f float64) bool {
	// is a natural number?
	if float64(int(num_f)) != num_f {
		return false
	}
	num := int(num_f)

	return big.NewInt(int64(num)).ProbablyPrime(0)

	// if num < 2 {
	// return false
	// }
	// if num == 2 {
	// return true
	// }
	// for x := 2; x < num/2; x += 2 {
	// if num%x == 0 {
	// return false
	// }
	// }
	// return true
}
