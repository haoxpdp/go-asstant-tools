package main

import "flag"

var i int

func main() {

	println(i)
}

func init() {
	flag.IntVar(&i, "i", 0, "i value")
	flag.Parse()
}
