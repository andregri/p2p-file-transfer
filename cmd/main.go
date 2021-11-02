package main

import (
	"context"
	"flag"
	"log"

	p2p "github.com/andregri/p2p-file-transfer/p2p"
	golog "github.com/ipfs/go-log/v2"
)

type Config struct {
	listenPort int
	peerId     string
	filePath   string
}

var logger = golog.Logger("file-transfer")

func main() {
	// Set context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set logger
	golog.SetLogLevel("file-transfer", "info")

	// Parse command-line
	var config Config
	flag.IntVar(&config.listenPort, "listen", 0, "port to listen incoming connections")
	flag.StringVar(&config.peerId, "peer", "", "peer we want to send the file")
	flag.StringVar(&config.filePath, "file", "", "path of the file to send")
	flag.Parse()

	if config.listenPort == 0 {
		log.Fatalln("Please, provide a valid port to bind on")
	}

	// Make new p2p host
	host, _ := p2p.MakeNewHost(ctx, config.listenPort)
	logger.Info("Host created. We are:", host.ID())
	logger.Info(host.Addrs())

	if config.peerId == "" {
		log.Println("sender")
	} else {
		log.Println("receiver")
	}
}
