package main

import (
	"log"
	"os"

	"github.com/mswift42/nip/cl"
)

func main() {
	app := cl.InitCli()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
