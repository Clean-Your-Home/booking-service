package main

import (
	"log"

	"booking-service/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("init: %v", err)
	}
	if err := a.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}
