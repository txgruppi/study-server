package cli

import (
	"fmt"
	"time"

	"github.com/txgruppi/study-server/http"
	"github.com/urfave/cli/v2"
)

func Wire(runner http.Runner) *cli.App {
	app := cli.NewApp()
	app.Flags = append(app.Flags, &cli.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Value:   8080,
		Usage:   "port to listen on",
		EnvVars: []string{"PORT"},
	})
	app.Flags = append(app.Flags, &cli.StringFlag{
		Name:    "delay-min",
		Value:   "0s",
		Usage:   "minimum delay for slow request simulation",
		EnvVars: []string{"DELAY_MIN"},
	})
	app.Flags = append(app.Flags, &cli.StringFlag{
		Name:    "delay-max",
		Value:   "5s",
		Usage:   "maximum delay for slow request simulation",
		EnvVars: []string{"DELAY_MAX"},
	})
	app.Action = func(c *cli.Context) error {
		port := c.Int("port")
		delayMinStr := c.String("delay-min")
		delayMaxStr := c.String("delay-max")
		delayMin, err := time.ParseDuration(delayMinStr)
		if err != nil {
			return fmt.Errorf("invalid delay-min: %w", err)
		}
		delayMax, err := time.ParseDuration(delayMaxStr)
		if err != nil {
			return fmt.Errorf("invalid delay-max: %w", err)
		}
		return runner(port, delayMin, delayMax)
	}
	return app
}
