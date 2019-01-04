package main

import "flag"

var verbose bool

func main() {
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	wristdevhack()
}
