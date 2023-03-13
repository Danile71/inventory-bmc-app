// Copyright 2023 NJWS Inc.

package discovery

import (
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/bootstrap"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/monitor"
	"github.com/urfave/cli/v2"
)

var (
	CLI = cli.NewApp()

	version = "v0.1.0"
)

func init() {
	CLI.Usage = "Discovery redfish tool"
	CLI.Version = version

	CLI.Commands = []*cli.Command{
		&cli.Command{
			Name:        "bootstrap",
			Description: "Bootstrap discovery",
			Action: func(ctx *cli.Context) (err error) {
				return bootstrap.Run(ctx.Context)
			},
		},
		&cli.Command{
			Name:        "monitor",
			Description: "Run discovery monitor",

			Action: func(ctx *cli.Context) (err error) {
				return monitor.Run(ctx.Context)
			},
		},
		&cli.Command{
			Name:        "run",
			Description: "Run discovery",

			Action: func(ctx *cli.Context) (err error) {
				if err = bootstrap.Run(ctx.Context); err != nil {
					return
				}
				return agent.Run(ctx.Context)
			},
		},
	}
}
