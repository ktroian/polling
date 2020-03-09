package main

import (
    "fmt"
    "log"
    "flag"
    "context"
    "os"
    "syscall"
    "os/signal"

    "github.com/ktroian/polling/daemon/action"
)

func handleSignals(cancel context.CancelFunc) {
	go func() {
		sigChannel := make(chan os.Signal)
		signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

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
    host 	 := flag.String("host", "localhost", "hostname for the service API")
    port 	 := flag.Int("port", 8080, "a port of the service API")
    interval := flag.Int("interval", 30, "an interval to check the API")
    uname 	 := flag.String("username", "postgres", "username to log into database")
    pswd 	 := flag.String("password", "postgres", "a password to log into database")
    logfile  := flag.String("log", "", "path to a file to write logs")

    flag.Parse()

	if len(*logfile) > 1 {
		file, err := os.OpenFile(*logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
		    log.Fatalf("error opening file: %v", err)
		}
		defer file.Close()

		log.SetOutput(file)
	}
    action.Init()

	ctx, cancel := context.WithCancel(context.Background())
    handleSignals(cancel)
    url := fmt.Sprintf("http://%s:%d/companies", *host, *port)
	connectToDatabase(ctx, *uname, *pswd)
    log.Println("Starting", os.Args[0])

    startPolling(ctx, url, *interval)
    defer log.Println(os.Args[0], "stopped")

	<- ctx.Done()
}