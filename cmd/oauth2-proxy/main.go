package main

import (
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/webhippie/oauth2-proxy/pkg/config"
	"github.com/webhippie/oauth2-proxy/pkg/version"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if env := os.Getenv("OAUTH2_PROXY_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	app := &cli.App{
		Name:     "oauth2-proxy",
		Version:  version.Version.String(),
		Usage:    "proxy for authentication via oauth2",
		Compiled: time.Now(),

		Authors: []*cli.Author{
			{
				Name:  "Thomas Boerger",
				Email: "thomas@webhippie.de",
			},
		},

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Value:       "info",
				Usage:       "set logging level",
				EnvVars:     []string{"KLEISTER_LOG_LEVEL"},
				Destination: &config.LogLevel,
			},
		},

		Before: func(c *cli.Context) error {
			return nil
		},

		Commands: []*cli.Command{
			Server(),
			Health(),
		},
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "show the help, so what you see now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print the current version of that tool",
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
