// Copyright 2023 NJWS Inc.

package agent

import (
	"time"

	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
)

func (a *Agent) entrypoint() (err error) {
	// wait router, need to register port
	time.Sleep(time.Millisecond * 50)

	if err = a.createOrUpdateFunctionLink(types.FunctionContainerPath, types.InventoryFunctionPath, types.InventoryFunctionLink); err != nil {
		return
	}

	if err = a.createOrUpdateFunctionLink(types.InventoryFunctionPath, types.ServiceFunctionPath, types.ServiceFunctionLink); err != nil {
		return
	}

	if err = a.createOrUpdateFunctionLink(types.InventoryFunctionPath, types.SystemsFunctionPath, types.SystemsFunctionLink); err != nil {
		return
	}

	if err = a.createOrUpdateFunctionLink(types.InventoryFunctionPath, types.SystemFunctionPath, types.SystemFunctionLink); err != nil {
		return
	}
	return
}
