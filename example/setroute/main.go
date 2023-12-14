package main

import (
	"log"
	mikrotik "mikrotik/mkrotik-api"
)

func main() {
	address := "1.1.1.1"
	username := "admin"
	password := "admin"

	// Create a new Mikrotik instance
	mikrotikClient, err := mikrotik.NewMikrotikRepository(address, username, password)
	if err != nil {
		log.Fatalf("Failed to connect to Mikrotik: %v", err)
	}
	err = mikrotikClient.SetStaticRoute("2.2.2.0/24", "blackhole", "1.1.1.1", "18013:13")
	if err != nil {
		log.Fatalf("Set route to Mikrotik: %v", err)
	}
}
