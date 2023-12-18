package main

import (
	"log"
	mikrotik "mikrotik/mkrotik-api"
)

func main() {
	address := "1.1.1.1"
	username := "admin"
	password := "admin"
	port := "8728"

	// Create a new Mikrotik instance
	mikrotikClient, err := mikrotik.NewMikrotikRepository(address, username, password, port)
	if err != nil {
		log.Fatalf("Failed to connect to Mikrotik: %v", err)
	}
	dstAddrs := []string{"2.2.2.0/24", "2.2.3.0/24", "1.2.2.0/24"}
	err = mikrotikClient.SetStaticRoute(dstAddrs, "no")
	if err != nil {
		log.Fatalf("Set route to Mikrotik: %v", err)
	}
}
