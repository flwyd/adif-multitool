// adifmt provides a variety of subcommands for manipulating ADIF log files.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("TODO: Implement adifmt commands")
	} else {
		cmd := os.Args[1]
		fmt.Printf("TODO: Implement adif %s\n", cmd)
	}
}
