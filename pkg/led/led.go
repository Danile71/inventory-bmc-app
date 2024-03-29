// Copyright 2023 NJWS Inc.

package led

import (
	"github.com/foliagecp/inventory-bmc-app/pkg/led/agent"
	"github.com/foliagecp/inventory-bmc-app/pkg/led/bootstrap"
	"github.com/urfave/cli/v2"
)

var (
	CLI = cli.NewApp()

	version = "v0.1.0"
)

func init() {
	CLI.Usage = "Update chassis indicator LED tool"
	CLI.Version = version
	CLI.Commands = []*cli.Command{
		&cli.Command{
			Name:        "bootstrap",
			Description: "Bootstrap update chassis indicator LED tool",
			Action: func(ctx *cli.Context) (err error) {
				return bootstrap.Run(ctx.Context)
			},
		},
		&cli.Command{
			Name:        "run",
			Description: "Run update chassis indicator LED tool",
			Action: func(ctx *cli.Context) (err error) {
				if err = bootstrap.Run(ctx.Context); err != nil {
					return
				}
				return agent.Run(ctx.Context)
			},
		},
	}
}
