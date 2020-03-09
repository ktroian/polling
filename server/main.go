package main

import (
	"os"
	"log"
	"flag"
	"context"
	"os/signal"
)

func handleSignals(cancel context.CancelFunc) {
	// to allow graceful shut downing
	go func() {
		sigChannel := make(chan os.Signal)
		signal.Notify(sigChannel, os.Interrupt)

		for {
			sig := <- sigChannel
			switch sig {
			case os.Interrupt:
				cancel()
				return
			}
		}
	}()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	updChannel := make(chan bool)
	companies := &AllCompanies{ companies: &Companies{} }

	file := flag.String("file", "", "path to a .csv file")
	port := flag.Int("port", 8080, "a port of the server")
	defer log.Println("Server Stopped")

	handleSignals(cancel)
	flag.Parse()

	initWriter(ctx, updChannel, companies, *file)
	serveRoutes(ctx, updChannel, companies, *port)

	<- ctx.Done()
}
