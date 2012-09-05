// LICENSE HERE

// proxy backwards and forwards
package main

import (
	"fmt"
	"flag"
	"net"
	"os"
	"bufio"
)

//----------------------------------------------------------------------
// linking functions
func pull_conn(conn net.Conn, c chan byte) {
	read := bufio.NewReader(conn)
	for {
		byt, err := read.ReadByte()
		if err != nil {
		}
		c <- byt
	}
}

func push_conn(conn net.Conn, c chan byte) {
	writer := bufio.NewWriter(conn)
	for {
		writer.WriteByte(<-c)
		writer.Flush()
	}
}

func linkcross() {
	c := make(chan byte)
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
	}
	conn, err := ln.Accept()
	if err != nil {
	}
	go push_conn(conn, c)
	chars :=  []byte("FUCK")
	for i := 0; i < len(chars); i++ {
		c <- chars[i]
	}
}

// func linkcross_onn(list_conn net.Conn, conn_addr string) {
// 	conn, err := net.Dial("tcp", conn_addr)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	fmt.Fprintf(list_conn, )
// }

// func linkcross(list_addr, conn_addr string) {
// 	ln, err := net.Listen("tcp", list_addr)
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	for {
// 		conn, err := ln.Accept()
// 		if err != nil {
// 			fmt.Print(err)
// 			continue
// 		}
// 		go linkcross_one(conn, conn_addr)
// 	}
// }

//----------------------------------------------------------------------
// get us some addresses
type addresses []string
func (addrs *addresses) String() string {
	return fmt.Sprint(*addrs)
}
func (addrs *addresses) Set(value string) error {
	*addrs = append(*addrs, value)
	return nil
}
var listenersFlag addresses
var connectorsFlag addresses

func main() {
	flag.Var(&listenersFlag, "l", "Which ports to listen on")
	flag.Var(&connectorsFlag, "c", "Which addresses to try to connect to")
	flag.Parse()

	if len(listenersFlag) + len(connectorsFlag) != 2 {
		fmt.Fprintln(os.Stderr, "Only 2 connections allowed")
		os.Exit(1)
	}
	fmt.Println(listenersFlag)
	fmt.Println(connectorsFlag)

	if len(listenersFlag) == 1 && len(connectorsFlag) == 1 {
		fmt.Println("FUCK")
	}
	if len(listenersFlag) == 40 {
		net.Dial("tcp", "google.com:8080")
	}
	linkcross()
}
