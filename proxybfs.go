// LICENSE HERE

// proxy backwards and forwards
package main

import (
	"fmt"
	"os"
)

// takes arguments, gives back a list of listeners, connectors, and
// dict of other args
func parseargs(args []string) (map[string]bool, map[string]bool, map[string]string) {
	// using maps as a set
	listeners := make(map[string]bool)
	connectors := make(map[string]bool)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-l":
			listeners[args[i+1]] = true
		case "-c":
			connectors[args[i+1]] = true
		default:
			fmt.Println("OH FUCK")
		}
	}
	kwargs := make(map[string]string)
	return listeners, connectors, kwargs
}

func main() {
	l,c,kw := parseargs(os.Args)
	fmt.Println(l,c,kw)
}
