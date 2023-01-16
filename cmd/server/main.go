package main

import (
	"context"
	"log"

	"github.com/walkergriggs/enoki/enoki/server"
)

func main() {
	ctx := context.Background()

	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
