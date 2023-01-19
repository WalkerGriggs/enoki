package main

import (
	"context"
	"log"

	"github.com/walkergriggs/enoki/internal/servers/gateway"
)

func main() {
	ctx := context.Background()

	if err := gateway.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
