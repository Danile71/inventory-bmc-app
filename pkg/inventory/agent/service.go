// Copyright 2023 NJWS Inc.

package agent

import (
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types/message"
)

const serviceMask = "service.*[?@._id == '%s'?].objects.root"

// serviceFunction flink function
func (a *Agent) serviceFunction(ctx module.Context) (err error) {
	client, err := a.getClientFromCmdb(ctx)
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument(serviceMask, ctx.Self().Id)
	if err != nil {
		if functionContext, err = system.CreateChild(ctx.Self().Id, types.RedfishServiceID, types.RedfishServiceLink, client.GetService()); err != nil {
			return err
		}
	} else {
		if functionContext, err = system.UpdateObject(document.Id.String(), client.GetService()); err != nil {
			return err
		}
	}

	// need to get execute 'onServiceFunctionResult'
	functionContext.ReplyResult = ctx.GetReplyResult()

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}
	
	client.GetService().Systems()

	ctx.Send(msg)
	return
}

// onServiceFunctionResult recv result after createOrUpdateService
func (a *Agent) onServiceFunctionResult(ctx module.Context, functionResult *pbtypes.FunctionResult) {
	if !functionResult.Complete {
		return
	}

	m, err := a.getMessageFromContext(ctx)
	if err != nil {
		return
	}

	if err := a.inventorySystems(ctx, m); err != nil {
		ctx.AddError(err)
	}
}

func (a *Agent) inventorySystems(ctx module.Context, m *message.Message) (err error) {
	document, err := a.getDocument(serviceMask, ctx.Self().Id)
	if err != nil {
		return
	}

	functionContext, err := prepareSystemsFunc(document.Id.String(), m)
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
