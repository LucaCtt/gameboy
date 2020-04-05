package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04"}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
}
