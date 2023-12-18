package main

import (
	"fmt"
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
	prefSrc := []string{"1.1.1.1", "1.1.1.2"}
	rlist := mikrotikClient.GetRouteList(prefSrc)
	fmt.Println(rlist)
}
