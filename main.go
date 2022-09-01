package main

import (
	"log"
	"os"

	"github.com/txgruppi/study-server/app"
)

func run() error {
	app, err := app.Wire()
	if err != nil {
		return err
	}
	return app.Run(os.Args)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
