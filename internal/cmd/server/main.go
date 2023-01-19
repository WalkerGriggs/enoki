package main

import (
	"context"
	"log"

	"github.com/walkergriggs/enoki/internal/servers/manifest"
)

func main() {
	ctx := context.Background()

	if err := manifest.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
