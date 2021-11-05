package main

import (
	"bufio"
	"context"
	"flag"
	"log"

	p2p "github.com/andregri/p2p-file-transfer/p2p"
	golog "github.com/ipfs/go-log/v2"
	net "github.com/libp2p/go-libp2p-core/network"
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
		log.Println("receiver")

		host.SetStreamHandler("/filetransfer/1.0.0", func(s net.Stream) {
			// Create a buffer stream for non blocking read and write
			rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

			go p2p.RecvFile(rw)
		})

	} else {
		log.Println("sender")

		id := p2p.Connect(ctx, host, config.peerId)

		// Open stream from this host to the target host
		s, err := host.NewStream(ctx, id, "/filetransfer/1.0.0")
		if err != nil {
			logger.Warn(err)
			return
		}
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
		go p2p.SendFile(rw, config.filePath)
	}

	select {}
}
