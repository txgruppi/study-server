package cli

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
)

func Wire(http *echo.Echo) *cli.App {
	app := cli.NewApp()
	app.Flags = append(app.Flags, &cli.IntFlag{
		Name:    "port",
		Aliases: []string{"p"},
		Value:   8080,
		Usage:   "port to listen on",
		EnvVars: []string{"PORT"},
	})
	app.Action = func(c *cli.Context) error {
		port := c.Int("port")
		return http.Start(fmt.Sprintf(":%d", port))
	}
	return app
}
