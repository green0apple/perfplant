package main

import (
	"fmt"
	"perfplant/perfs/udp"
)

func main() {
	c := udp.NewClient()
	fmt.Printf("%s\n", c.Run())
}
