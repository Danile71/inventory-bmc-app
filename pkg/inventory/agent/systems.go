// Copyright 2023 NJWS Inc.

package agent

import (
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/stmcginnis/gofish/common"
)

const systemsCollection = "/redfish/v1/Systems"

// systemsFunction flink function
func (a *Agent) systemsFunction(ctx module.Context) (err error) {
	client, err := a.getClientFromContext(ctx)
	if err != nil {
		return
	}

	links, err := common.GetCollection(client, systemsCollection)
	if err != nil {
		return
	}

	m, err := a.getMessageFromContext(ctx)
	if err != nil {
		return
	}

	for _, link := range links.ItemLinks {
		m.Additional = link

		functionContext, err := prepareSystemFunc(ctx.Self().Id, m)
		if err != nil {
			return err
		}

		msg, err := module.ToMessage(functionContext)
		if err != nil {
			return err
		}
		ctx.Send(msg)
	}

	return
}
