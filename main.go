package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

const DefentTimeOut = time.Second * 10

var UnsupportedProtocols = errors.New("Unsupported protocols")

func printErr(us string, e error) {
	fmt.Fprintln(os.Stderr, e, ":", us)
}

func errHandle(us string, e error) bool {
	if e != nil {
		printErr(us, e)
		return true
	}
	return false
}

func httpHandle(us string, u *url.URL) {
	client := &http.Client{
		Timeout: DefentTimeOut,
	}
	resp, err := client.Get(us)
	if errHandle(us, err) {
		return
	}
	defer resp.Body.Close()
	fmt.Fprintln(os.Stdout, us)
}

func tcpHandle(us string, u *url.URL) {
	conn, err := net.DialTimeout("tcp", u.Host, DefentTimeOut)
	if errHandle(us, err) {
		return
	}
	defer conn.Close()
	fmt.Fprintln(os.Stdout, us)
}

func udpHandle(us string, u *url.URL) {
	conn, err := net.DialTimeout("udp", u.Host, DefentTimeOut)
	if errHandle(us, err) {
		return
	}
	defer conn.Close()
	fmt.Fprintln(os.Stdout, us)
}

func wsHandle(us string, u *url.URL) {
	ctx, cancel := context.WithTimeout(context.Background(), DefentTimeOut)
	defer cancel()

	c, _, err := websocket.Dial(ctx, us, nil)
	if errHandle(us, err) {
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")
	fmt.Fprintln(os.Stdout, us)
}

func scan(us string, wg *sync.WaitGroup, ch chan struct{}) {
	defer wg.Done()
	u, err := url.Parse(us)
	if errHandle(us, err) {
		return
	}
	switch u.Scheme {
	case "http":
		httpHandle(us, u)
	case "https":
		httpHandle(us, u)
	case "tcp":
		tcpHandle(us, u)
	case "udp":
		udpHandle(us, u)
	case "ws":
		wsHandle(us, u)
	case "wss":
		wsHandle(us, u)
	default:
		printErr(us, UnsupportedProtocols)
	}
	<-ch
}

func main() {

	arg := new(struct {
		Help              bool
		ConcurrencyLimits int
		File              string
		Out               string
		ErrOut            string
	})

	flag.BoolVar(&arg.Help, "h", false, "this help.")
	flag.BoolVar(&arg.Help, "help", false, "this help.")
	flag.IntVar(&arg.ConcurrencyLimits, "c", 1, "Concurrency limits.")
	flag.Parse()

	if arg.Help {
		flag.Usage()
		return
	}
	if arg.ConcurrencyLimits < 1 {
		fmt.Fprintln(os.Stderr, "Concurrency limits must be greater than 0.")
		return
	}
	var us string
	limitChan := make(chan struct{}, arg.ConcurrencyLimits)
	wg := sync.WaitGroup{}
	for {
		_, err := fmt.Scan(&us)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		limitChan <- struct{}{}
		wg.Add(1)
		go scan(us, &wg, limitChan)
	}
	wg.Wait()
}
