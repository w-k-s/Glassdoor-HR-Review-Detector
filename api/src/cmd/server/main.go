package main

import (
	"log"

	s "com.github/w-k-s/glassdoor-hr-review-detector/internal/server"
)

func main() {
	config := s.ReadConfig()
	server := s.NewServer(config)
	log.Printf("Server running on port %q", config.Server.ListenAddress)
	log.Fatal(server.Start())
}
