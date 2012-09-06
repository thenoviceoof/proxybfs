// LICENSE HERE

// proxy backwards and forwards
package main

import (
	"fmt"
	"flag"
	"net"
	"os"
	"io"  // for the EOF error
	"bufio"
)

//----------------------------------------------------------------------
// linking functions
func pull_conn(conn net.Conn, c chan byte, closed chan bool) {
	read := bufio.NewReader(conn)
	for {
		byt, err := read.ReadByte()
		switch err {
		case io.EOF:
			break
		default:
			fmt.Println(err)
		}
		fmt.Print("got byte: " + string(byt))
		c <- byt
	}
	close(c)
	closed <- true
}

func push_conn(conn net.Conn, c chan byte) {
	writer := bufio.NewWriter(conn)
	for {
		byt, ok := <-c
		if !ok {
			break
		}
		fmt.Print("putting byte: " + string(byt))
		writer.WriteByte(byt)
		writer.Flush()
	}
	conn.Close()
}

// facilitate trading between two connections
func crosspipe(pipea, pipeb net.Conn) {
	fmt.Println("Linking up")
	a2b := make(chan byte)
	b2a := make(chan byte)
	finish := make(chan bool)
	go pull_conn(pipea, a2b, finish)
	go pull_conn(pipeb, b2a, finish)
	go push_conn(pipea, b2a)
	go push_conn(pipeb, a2b)
	fmt.Println("Finishing up...")
	_, _ = <-finish, <-finish
	fmt.Println("And there you have it!")
}

// for each listening connection, make an outgoing one
func listenTo(list_addr, conn_addr string) {
	ln, err := net.Listen("tcp", list_addr)
	if err != nil {
		fmt.Println(err)
	}
	// keep accepting connections
	for {
		ln_conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("dialing")
		cn_conn, err := net.Dial("tcp", conn_addr)
		fmt.Println("done dialing")
		go crosspipe(ln_conn, cn_conn)
		fmt.Println("waiting for new conn...")
	}
}

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
		listenTo(listenersFlag[0], connectorsFlag[0])
	}
	if len(listenersFlag) == 40 {
		net.Dial("tcp", "google.com:8080")
	}
}
