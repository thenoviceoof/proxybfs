// LICENSE HERE

// proxy backwards and forwards
package main

import (
	"fmt"
	"flag"
	"net"
	"os"
)

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
		fmt.Fprint(os.Stderr, "Too many connections")
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
}
