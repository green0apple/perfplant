package perf

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"perfplant/perf/udp"
	"syscall"
	"time"
)

const (
	PROTOCOL_TYPE_UDP   = "udp"
	PROTOCOL_TYPE_TCP   = "tcp"
	PROTOCOL_TYPE_TLS   = "tls"
	PROTOCOL_TYPE_HTTP  = "http"
	PROTOCOL_TYPE_HTTPS = "https"
	PROTOCOL_TYPE_H2    = "h2"
	PROTOOCL_TYPE_H3    = "h3"
)

var (
	ErrUnknownProtocolType error = errors.New("unknown protocol type")
	ErrPlantAlreadyRunning error = errors.New("plant already running")
)

type Client interface {
	Request()
	Response()
	RequestedCount()
}

type Plant struct {
	ProtocolType        string
	RPS                 int64
	MaxRequstPerSession int
	MaxWorkers          int

	order              chan struct{}
	stopRequestControl chan struct{}
	stop               chan struct{}
	running            bool
}

func (p *Plant) Run() error {
	// TODO :: return report

	if p.running {
		return ErrPlantAlreadyRunning
	}
	p.running = true
	defer p.doStop()

	switch p.ProtocolType {
	default:
		return ErrUnknownProtocolType
	case PROTOCOL_TYPE_UDP:
	case PROTOCOL_TYPE_TCP:
	case PROTOCOL_TYPE_TLS:
	case PROTOCOL_TYPE_HTTP:
	case PROTOCOL_TYPE_HTTPS:
	case PROTOCOL_TYPE_H2:
	case PROTOOCL_TYPE_H3:
	}

	p.order = make(chan struct{})
	p.stopRequestControl = make(chan struct{}, 1)
	go p.requestControl()

	stopWorkers := make([]chan struct{}, p.MaxWorkers)

	for i := 0; i < p.MaxWorkers; i++ {
		stopWorkers[i] = make(chan struct{}, 1)
		go func(stop chan struct{}) {
			for {
				select {
				case <-stop:
					return
				case <-p.order:
					go p.request()
				}
			}
		}(stopWorkers[i])
	}
	defer func() {
		for i := 0; i < p.MaxWorkers; i++ {
			stopWorkers[i] <- struct{}{}
			close(stopWorkers[i])
		}
	}()

	stopBySignal := make(chan os.Signal, 1)
	signal.Notify(stopBySignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT)

	p.stop = make(chan struct{}, 1)
	for {
		select {
		case <-p.stop:
			fmt.Printf("stopped by calling Stop()\n")
			return nil
		case sig := <-stopBySignal:
			close(p.stop)
			p.stop = nil
			fmt.Printf("stopped by signal %s\n", sig.String())
			return nil
		}
	}
}

func (p *Plant) Stop() {
	if p.running {
		p.doStop()
	}
}

func (p *Plant) doStop() {
	if p.stopRequestControl != nil {
		p.stopRequestControl <- struct{}{}
		close(p.stopRequestControl)
		p.stopRequestControl = nil
	}

	if p.order != nil {
		close(p.order)
		p.order = nil
	}

	if p.stop != nil {
		close(p.stop)
		p.stop = nil
	}
	p.running = false
}

// requestControl() controls request per second
func (p *Plant) requestControl() {
	var (
		expectedIntervalNano time.Duration = time.Duration(1000000000 / p.RPS)
		totalShot            int64         = 0
		timeToWait           time.Duration
		ticker               *time.Ticker = time.NewTicker(1)

		expectedShot int64
		startedAt    int64
	)

	defer ticker.Stop()

	startedAt = time.Now().UnixNano()

	fmt.Println("start at", time.Now().String())
	for {
		expectedShot = (p.RPS * (time.Now().UnixNano() - startedAt)) / 1000000000

		if totalShot < expectedShot {
			timeToWait = 1
		} else {
			timeToWait = time.Duration(expectedIntervalNano)
		}

		ticker.Reset(timeToWait)
		select {
		case <-p.stopRequestControl:
			return
		case <-ticker.C:
			p.order <- struct{}{}
			totalShot++
		}

		if totalShot == p.RPS*60 {
			fmt.Println("done at", time.Now().String())
			break
		}
	}
}

func (p *Plant) request() {
	 := udp.NewClient()
	c.Dial()
	c.Write()
}
