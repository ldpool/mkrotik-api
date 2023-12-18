package main

import (
	"log"
	mikrotik "mikrotik/mkrotik-api"
)

func main() {
	address := "122.226.180.206"
	username := "admin"
	password := "mmdoudou"
	port := "18728"

	// Create a new Mikrotik instance
	mikrotikClient, err := mikrotik.NewMikrotikRepository(address, username, password, port)
	if err != nil {
		log.Fatalf("Failed to connect to Mikrotik: %v", err)
	}
	err = mikrotikClient.AddStaticRoute("1.2.2.0/24", "unreachable", "1.1.1.1", "18013:12,18013:13", "cloud")
	if err != nil {
		log.Fatalf("Set route to Mikrotik: %v", err)
	}
}
