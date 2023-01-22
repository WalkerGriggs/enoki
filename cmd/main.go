package main

import (
	"context"
	"log"

	"github.com/walkergriggs/enoki/internal/servers/gateway"
	"github.com/walkergriggs/enoki/internal/servers/manifest"
	"github.com/walkergriggs/enoki/internal/servers/storage"
)

func main() {
	ctx := context.Background()

	log.Println("Booting up")

	go func() {
		if err := manifest.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := storage.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	
	if err := gateway.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
