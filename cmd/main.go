package main

import (
	"log"
	uf "unri_fusioner"
)

func main() {
	configOpt := func(c *uf.Config) {
		c.SintaDomain = "https://sinta.kemdikbud.go.id"
		c.Host = "127.0.0.1"
		c.Port = 3000
	}
	config := uf.LoadConfig(configOpt)

	httpServerOpt := func(s *HTTPServer) error {
		s.host = config.Host
		s.port = config.Port
		return nil
	}

	server, err := NewHTTPServer(config, httpServerOpt)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server running on %s port %d ...", server.host, server.port)
	server.Start()
}
