package main

import (
	"context"
	"log"
	
	"github.com/walkergriggs/enoki/enoki/gateway"
)

func main() {
	ctx := context.Background()

	if err := gateway.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
