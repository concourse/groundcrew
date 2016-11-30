package main

import (
	"os"

	"github.com/concourse/groundcrew/commands"
	"github.com/jessevdk/go-flags"
)

func main() {
	cmd := &commands.GroundcrewCommand{}

	parser := flags.NewParser(cmd, flags.Default)
	parser.NamespaceDelimiter = "-"

	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}
}
