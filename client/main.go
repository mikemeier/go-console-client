package main

import (
	"console/client/app"
	"console/client/socket"
	"flag"
)

var local = flag.Bool("dev", false, "Connect client to localhost")

func main() {
	flag.Parse()

	var host = "thetrader.ch:9191"
	if *local {
		host = "localhost:9191"
	}

	done := make(chan bool)

	go app.Run(host)
	go socket.Heartbeat(host, !*local)
	go socket.ListenForData()

	done <- false
}
