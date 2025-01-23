package main

import (
	"fmt"
	"log"

	"github.com/MatthewAraujo/airCast/cmd/api"
)

func main() {
	server := api.NewAPIServer(fmt.Sprintf(":%d", 8080))
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
