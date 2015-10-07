package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hyqhyq3/through"
	"github.com/hyqhyq3/through/proxy/http"
	"github.com/hyqhyq3/through/proxy/socks5"
)

func serveSocks5(addr string, exit <-chan bool) {
	if err := socks5_proxy.ListenAndServeSocks5(addr, &through.RouteDialer{}, exit); err != nil {
		log.Fatal(err)
	}
}

func serveHTTP(addr string, exit <-chan bool) {
	if err := http_proxy.ListenAndServe(addr, &through.RouteDialer{}, exit); err != nil {
		log.Fatal(err)
	}
}

var http_addr, socks5_addr string

func startServer(exit <-chan bool) {

	if len(http_addr) > 0 {
		go serveHTTP(http_addr, exit)
	}
	if len(socks5_addr) > 0 {
		go serveSocks5(socks5_addr, exit)
	}

}

func init() {
	flag.StringVar(&http_addr, "http", "", "bind http proxy")
	flag.StringVar(&socks5_addr, "socks5", "", "bind socks5 proxy")
	flag.Parse()
	fmt.Println(http_addr, socks5_addr)
}

func main() {

	through.InitConfig("config.ini")

	exit := make(chan bool, 20)
	startServer(exit)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Println("received ", sig.String())
	exit <- true
	close(exit)
}
