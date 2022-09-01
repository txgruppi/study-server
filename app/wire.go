//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/txgruppi/study-server/cli"
	"github.com/txgruppi/study-server/database"
	"github.com/txgruppi/study-server/http"
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
