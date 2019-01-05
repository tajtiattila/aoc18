package main

import "flag"

var verbose bool

func main() {
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.Parse()

	wristdevhack()      // 16
	reservoirresearch() // 17
	collectlumber()     // 18
}
