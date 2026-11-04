package main

import (
	"log"

	"github.com/firefly-zero/fonts.fireflyzero.com/generator"
)

func main() {
	err := generator.Run()
	if err != nil {
		log.Fatalf("fatal error: %v", err)
	}
}
