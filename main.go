package main

import (
	"perfplant/perf"
)

func main() {
	p := perf.Plant{
		ProtocolType: perf.PROTOCOL_TYPE_UDP,
		RPS:          1,
		MaxWorkers:   1,
	}

	p.Run()
}
