// LICENSE HERE

// proxy backwards and forwards
package main

import (
	"fmt"
	"os"
)

// takes arguments, gives back a list of listeners, connectors, and
// dict of other args
func parseargs(args []string) ([]string, []string, map[string]string) {
	listeners := []string{"world"}
	connectors := []string{"hello"}
	kwargs := make(map[string]string)
	return listeners, connectors, kwargs
}

func main() {
	parseargs(os.Args)
	fmt.Println("Falling over my own feet here")
}