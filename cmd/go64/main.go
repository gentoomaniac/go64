package main

import (
	"github.com/alecthomas/kong"
	"github.com/rs/zerolog/log"

	"github.com/gentoomaniac/gocli"
	"github.com/gentoomaniac/logging"

	"github.com/gentoomaniac/go64/pkg/c64"
)

var (
	version = "unknown"
	commit  = "unknown"
	binName = "unknown"
	builtBy = "unknown"
	date    = "unknown"
)

var cli struct {
	logging.LoggingConfig

	Run struct {
		BasicRom     string `help:"Path to the BASIC ROM" type:"existingfile" required:""`
		KernalRom    string `help:"Path to the Kernal ROM" type:"existingfile" required:""`
		CharacterRom string `help:"Path to the character ROM" type:"existingfile" required:""`
	} `cmd:"" help:"Run the application (default)." default:"1" hidden:""`

	Version gocli.VersionFlag `short:"V" help:"Display version."`
}

func main() {
	ctx := kong.Parse(&cli, kong.UsageOnError(), kong.Vars{
		"version": version,
		"commit":  commit,
		"binName": binName,
		"builtBy": builtBy,
		"date":    date,
	})
	logging.Setup(&cli.LoggingConfig)

	switch ctx.Command() {
	case "foo":
		log.Info().Msg("foo command")
	default:
		system := &c64.C64{}

		system.Init(cli.Run.BasicRom, cli.Run.BasicRom, cli.Run.CharacterRom)
		system.Run()
	}
	ctx.Exit(0)
}
