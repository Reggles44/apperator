package main

import (
	"log"
	"os"

	"github.com/Reggles44/apperator/pkg/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
