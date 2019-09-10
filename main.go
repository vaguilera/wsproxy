package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func customUsage() {
	binfile := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s [OPTIONS]\n", binfile)
	fmt.Printf("Ex: %s -p 8080 -r myserver.com -P 23000\n", binfile)
	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = customUsage
	log.SetFlags(0)
	localHost := flag.String("a", "localhost", "Listen address")
	localPort := flag.Int("p", 8080, "Listen port")
	remoteHost := flag.String("r", "", "TCP server address")
	remotePort := flag.Int("P", 21000, "TCP server Port")
	textMode := flag.Bool("t", true, "Text encoding (binary/text)")
	showHelp := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *remoteHost == "" {
		fmt.Printf("Error: you should specify a remote address\n")
		flag.Usage()
		os.Exit(1)
	}

	ListenWebSocket(*localHost, *localPort, *remoteHost, *remotePort, *textMode)
}
