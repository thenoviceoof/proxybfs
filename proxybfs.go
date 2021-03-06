/* MIT LICENSE
 * Copyright (c) <2012> <thenoviceoof>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation
 * files (the "Software"), to deal in the Software without
 * restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies
 * of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// proxy backwards and forwards
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"time"
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
func errmsg(code int, msg string) {
	fmt.Fprintln(os.Stderr, "[ERROR] "+msg)
	os.Exit(code)
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
	finish := make(chan bool) // tell this fn we're done
	go pull_conn(pipea, a2b, finish)
	go pull_conn(pipeb, b2a, finish)
	go push_conn(pipea, b2a)
	go push_conn(pipeb, a2b)
	_ = <-finish
}

// for each listening connection, make an outgoing one
func listenOne(list_addr, conn_addr string) {
	info("Starting to listen on: " + list_addr)
	ln, err := net.Listen("tcp", list_addr)
	if err != nil {
		errmsg(10, "Can't start server: "+err.Error())
	}
	// keep accepting connections
	for {
		ln_conn, err := ln.Accept()
		if err != nil {
			info(err.Error())
			continue
		}
		debug("dialing")
		cn_conn, err := net.Dial("tcp", conn_addr)
		if err != nil {
			info(err.Error())
			ln_conn.Close()
			continue
		}
		info("Connection created with " + conn_addr)
		go crosspipe(ln_conn, cn_conn)
		debug("waiting for new connection...")
	}
}

// listen on both sides, link them up
func listenTwo(lista_addr, listb_addr string) {
	info("Starting to listen on: " + lista_addr)
	lna, err := net.Listen("tcp", lista_addr)
	if err != nil {
		errmsg(10, "Can't start server: "+err.Error())
	}
	info("Starting to listen on: " + listb_addr)
	lnb, err := net.Listen("tcp", listb_addr)
	if err != nil {
		errmsg(10, "Can't start server: "+err.Error())
	}
	// keep accepting connections
	for {
		debug("waiting for new connection...")
		lna_conn, err := lna.Accept()
		if err != nil {
			info(err.Error())
			continue
		}
		lnb_conn, err := lnb.Accept()
		if err != nil {
			info(err.Error())
			lna.Close()
			continue
		}
		go crosspipe(lna_conn, lnb_conn)
	}
}

// time between failed connection tries
var retryPeriod time.Duration

// time between successful connection tries
var connPeriod time.Duration

// serially connect to two addresses
func connectTwo(conna_addr, connb_addr string) {
	for {
		info("Attempting to dial both...")
		conna, err := net.Dial("tcp", conna_addr)
		if err != nil {
			info(err.Error())
			time.Sleep(retryPeriod * time.Millisecond)
			continue
		}
		connb, err := net.Dial("tcp", connb_addr)
		if err != nil {
			info(err.Error())
			conna.Close()
			time.Sleep(retryPeriod * time.Millisecond)
			continue
		}
		debug("Have both connections, ")
		// no go, need to be serial
		crosspipe(conna, connb)
		debug("Both connections closed")
		// wait for a shorter time
		time.Sleep(connPeriod * time.Millisecond)
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

func normalizeAddr(addr string) string {
	matchp, err := regexp.Match("^\\d{1,5}$", []byte(addr))
	if err != nil {
		errmsg(2, "Failure while trying to normalize the address")
	}
	if matchp {
		return ":" + addr
	}
	return addr
}

func main() {
	// do the actual parsing
	flag.Var(&listenersFlag, "l", "Which ports to listen on")
	flag.Var(&connectorsFlag, "c", "Which addresses to try to connect to")
	flag.BoolVar(&infoptr, "v", false, "Turn on verbose mode")
	flag.BoolVar(&debugptr, "vv", false, "Turn on extra verbose mode")
	retryPeriod = time.Duration(1000 * (*flag.Float64("rp", 5.0,
		"Retry rate for double connections")))
	connPeriod = time.Duration(1000 * (*flag.Float64("cp", 0.5,
		"Retry rate for double connections, on success")))
	flag.Parse()

	debug("Number of listeners: " + fmt.Sprint(len(listenersFlag)))
	debug("Number of connectors: " + fmt.Sprint(len(connectorsFlag)))
	// check a possibly temporary condition
	if len(listenersFlag)+len(connectorsFlag) != 2 {
		errmsg(1, "Strictly 2 connections allowed")
	}

	if len(listenersFlag) == 1 && len(connectorsFlag) == 1 {
		listenOne(normalizeAddr(listenersFlag[0]),
			normalizeAddr(connectorsFlag[0]))
	}
	if len(listenersFlag) == 2 && len(connectorsFlag) == 0 {
		listenTwo(normalizeAddr(listenersFlag[0]),
			normalizeAddr(listenersFlag[1]))
	}
	if len(listenersFlag) == 0 && len(connectorsFlag) == 2 {
		connectTwo(normalizeAddr(connectorsFlag[0]),
			normalizeAddr(connectorsFlag[1]))
	}
}
