package main

import (
	"fmt"
	"perfplant/perf/udp"
)

func main() {
	c := udp.NewClient()
	fmt.Printf("%s\n", c.Run())
}
