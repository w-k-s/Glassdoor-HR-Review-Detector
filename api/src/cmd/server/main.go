package main

import (
	"flag"
	"log"

	s "com.github/w-k-s/glassdoor-hr-review-detector/internal/server"
)


func main() {
	listenAddress := flag.String("listenAddress",":3000", "The host and port on which the sever will listen for requests e.g. localhost:3000")
	flag.Parse()

	server := s.NewServer(*listenAddress)
	log.Printf("Server running on port %q", *listenAddress)
	log.Fatal(server.Start())
}
