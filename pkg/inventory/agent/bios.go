// Copyright 2023 NJWS Inc.

package agent

import (
	"path"

	"git.fg-tech.ru/listware/go-core/pkg/module"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/stmcginnis/gofish/redfish"
)

// flink function
// systemsFunction (cmdbContext -> "%s.redfish-devices.root")
func (a *Agent) biosFunction(ctx module.Context) (err error) {
	client, err := a.getClientFromContext(ctx)
	if err != nil {
		return
	}

	msg, err := a.getMessageFromContext(ctx)
	if err != nil {
		return
	}

	bios, err := redfish.GetBios(client, path.Join(msg.Additional, "Bios"))
	if err != nil {
		return
	}
	_ = bios.Attributes
	return
}
func (a *Agent) onBiosFunctionResult(ctx module.Context, functionResult *pbtypes.FunctionResult) {
	if !functionResult.Complete {
		return
	}

}
