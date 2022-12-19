package core

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

func respondPing(command_arr []interface{}) ([]byte, error) {
	if len(command_arr) == 0 {
		return []byte("+PONG\r\n"), nil
	} else if len(command_arr) >= 1 {
		return []byte("-ERR wrong number of arguments for 'ping' command\r\n"), nil
	} else {
		arg := command_arr[0].(string)
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)), nil
	}
}

func readCommand(c net.Conn) (string, error) {
	// TODO: Max read in one shot is 512 bytes
	// To allow input > 512 bytes, then repeated read until
	// we get EOF or designated delimiter
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {

	val, _ := Decode([]byte(cmd))

	arr_val := val.([]interface{})
	command := strings.ToUpper(arr_val[0].(string))

	var response = []byte("$-1\r\n")
	// var err = nil

	if command == "PING" {
		// fmt.Println("YES")
		response, _ = respondPing(arr_val[1:])
	} else {
		// fmt.Println("NO")
	}

	if _, err := c.Write(response); err != nil {
		return err
	}
	return nil
}

func RunServer() {
	var Host string = "0.0.0.0"
	var Port int = 6379

	log.Println("Starting the 🐭 server on", Host, Port)

	var con_clients int = 0

	// listening to the configured host:port
	lsnr, err := net.Listen("tcp", Host+":"+strconv.Itoa(Port))
	if err != nil {
		panic(err)
	}

	for {
		// blocking call: waiting for the new client to connect
		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}

		// increment the number of concurrent clients
		con_clients += 1
		log.Println("client connected with address:", c.RemoteAddr(), "concurrent clients", con_clients)

		for {
			// over the socket, continuously read the command and print it out
			cmd, err := readCommand(c)
			if err != nil {
				c.Close()
				con_clients -= 1
				log.Println("client disconnected", c.RemoteAddr(), "concurrent clients", con_clients)
				if err == io.EOF {
					break
				}
				log.Println("err", err)
			}
			log.Println("command", cmd)
			if err = respond(cmd, c); err != nil {
				log.Print("err write:", err)
			}
		}
	}
}
