package main

import (
	"flag"
	"log"
	"miba/internal/api"
)

func main() {
	password := flag.String("password", "", "Password for login")

	flag.Parse()

	if *password == "" {
		log.Fatalln("Password cannot be empty")
	}

	c, err := api.NewAPIClient()
	if err != nil {
		log.Fatalln("New API Client error:", err)
	}

	if err := c.PingRouter(); err != nil {
		log.Fatalln("Ping error:", err)
	}

	if err := c.UpdateInformation(); err != nil {
		log.Fatalln("Update info error:", err)
	}

	stok, err := c.Login(*password)
	if err != nil {
		log.Fatalln("Login error:", err)
	}

	log.Println(stok)
}
