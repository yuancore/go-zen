package main

import (
	"log"

	"github.com/yuancore/go-zen/examples/quickstart/boot/serve"
)

func main() {
	if err := serve.Run(); err != nil {
		log.Fatalf("quickstart run failed: %v", err)
	}
}
