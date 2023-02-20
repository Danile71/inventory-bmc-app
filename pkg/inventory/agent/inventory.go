// Copyright 2023 NJWS Inc.

package agent

import (
	"git.fg-tech.ru/listware/go-core/pkg/module"
)

// inventoryFunction flink function
// create/update link from  "service.inventory.redfish.functions.root" and execute it
func (a *Agent) inventoryFunction(ctx module.Context) (err error) {
	functionContext, err := prepareServiceFunc(ctx.Self().Id)
	if err != nil {
		return
	}
	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}
	ctx.Send(msg)
	return
}
