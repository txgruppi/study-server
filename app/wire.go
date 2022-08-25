//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/txgruppi/tasks-server/cli"
	"github.com/txgruppi/tasks-server/database"
	"github.com/txgruppi/tasks-server/http"
	ucli "github.com/urfave/cli/v2"
)

func Wire() (*ucli.App, error) {
	wire.Build(
		cli.Wire,
		http.Wire,
		database.Wire,
	)
	return &ucli.App{}, nil
}
