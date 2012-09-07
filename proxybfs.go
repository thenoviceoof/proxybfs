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
// utils

var infoptr = false
var debugptr = false

func info(msg string) {
	if infoptr || debugptr {
		fmt.Println("[INFO] " + msg)
	}
}
func debug(msg string) {
	if debugptr {
		fmt.Println("[DEBUG] " + msg)
	}
}

//----------------------------------------------------------------------
// linking functions
func pull_conn(conn net.Conn, c chan byte, closed chan bool) {
	debug("Pulling from conn: " + conn.RemoteAddr().Network())
	read := bufio.NewReader(conn)
	for {
		byt, err := read.ReadByte()
		if err != nil {
			debug("IO Error: " + err.Error())
			info("Closing connection...")
			break
		}
		debug("got byte: " + string(byt))
		c <- byt
	}
	close(c)
	closed <- true
}

func push_conn(conn net.Conn, c chan byte) {
	debug("Pushing to conn: " + conn.RemoteAddr().Network())
	writer := bufio.NewWriter(conn)
	for byt := range c {
		debug("putting byte: " + string(byt))
		writer.WriteByte(byt)
		writer.Flush()
	}
	conn.Close()
}

// facilitate trading between two connections
func crosspipe(pipea, pipeb net.Conn) {
	debug("Linking up two net connections")
	a2b := make(chan byte)
	b2a := make(chan byte)
	finish := make(chan bool)  // tell this fn we're done
	go pull_conn(pipea, a2b, finish)
	go pull_conn(pipeb, b2a, finish)
	go push_conn(pipea, b2a)
	go push_conn(pipeb, a2b)
	_ = <-finish
}

// for each listening connection, make an outgoing one
func listenOne(list_addr, conn_addr string) {
	info("Starting to listen on: " + listenersFlag[0])
	ln, err := net.Listen("tcp", list_addr)
	if err != nil {
		debug(err.Error())
		return
	}
	// keep accepting connections
	for {
		ln_conn, err := ln.Accept()
		if err != nil {
			info(err.Error())
		}
		debug("dialing")
		cn_conn, err := net.Dial("tcp", conn_addr)
		info("Connection created with " + conn_addr)
		go crosspipe(ln_conn, cn_conn)
		debug("waiting for new conn...")
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
	// do the actual parsing
	flag.Var(&listenersFlag, "l", "Which ports to listen on")
	flag.Var(&connectorsFlag, "c", "Which addresses to try to connect to")
	flag.BoolVar(&infoptr, "v", false, "Turn on verbose mode")
	flag.BoolVar(&debugptr, "vv", false, "Turn on extra verbose mode")
	flag.Parse()

	debug("Number of listeners: " + string(len(listenersFlag)))
	debug("Number of connectors: " + string(len(connectorsFlag)))
	// check a possibly temporary condition
	if len(listenersFlag) + len(connectorsFlag) != 2 {
		fmt.Fprintln(os.Stderr, "Only 2 connections allowed")
		os.Exit(1)
	}

	if len(listenersFlag) == 1 && len(connectorsFlag) == 1 {
		listenOne(listenersFlag[0], connectorsFlag[0])
	}
	if len(listenersFlag) == 40 {
		net.Dial("tcp", "google.com:8080")
	}
}
