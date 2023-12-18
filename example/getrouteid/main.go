package main

import (
	"fmt"
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
	dstAddrs := []string{"2.2.2.0/24", "2.2.3.0/24"}
	rId := mikrotikClient.GetRoutes(dstAddrs)
	fmt.Println(rId)
}
