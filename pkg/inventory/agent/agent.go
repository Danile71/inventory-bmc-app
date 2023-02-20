// Copyright 2023 NJWS Inc.

package agent

import (
	"context"
	"strings"

	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/apache/flink-statefun/statefun-sdk-go/v3/pkg/statefun"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

type Agent struct {
	ctx    context.Context
	cancel context.CancelFunc

	executor executor.Executor

	m module.Module

	sessionSpec statefun.ValueSpec

	genericSpec statefun.ValueSpec
}

// Run agent
func Run() (err error) {
	a := &Agent{
		sessionSpec: statefun.ValueSpec{
			Name:      "inventory_session_json",
			ValueType: module.GenericJsonType,
		},

		genericSpec: statefun.ValueSpec{
			Name:      "inventory_generic_json",
			ValueType: module.GenericJsonType,
		},
	}
	a.ctx, a.cancel = context.WithCancel(context.Background())

	if a.executor, err = executor.New(); err != nil {
		return
	}

	return a.run()
}

func appendPath(paths ...string) string {
	return strings.Join(paths, ".")
}

func (a *Agent) run() (err error) {
	defer a.executor.Close()

	log.Infof("run system agent")

	a.osSignalCtrl()

	a.m = module.New(types.Namespace, module.WithPort(31001))

	log.Infof("use port (%d)", a.m.Port())

	if err = a.m.Bind(types.InventoryFunctionType, a.inventoryFunction); err != nil {
		return
	}

	if err = a.m.Bind(types.ServiceFunctionType, a.serviceFunction, module.WithOnResult(a.onServiceFunctionResult), module.WithValueSpec(a.sessionSpec)); err != nil {
		return
	}

	if err = a.m.Bind(types.SystemsFunctionType, a.systemsFunction, module.WithValueSpec(a.sessionSpec)); err != nil {
		return
	}

	if err = a.m.Bind(types.SystemFunctionType, a.systemFunction, module.WithOnResult(a.onSystemFunctionResult), module.WithValueSpec(a.sessionSpec), module.WithValueSpec(a.genericSpec)); err != nil {
		return
	}

	if err = a.m.Bind(types.BiosFunctionType, a.biosFunction, module.WithOnResult(a.onBiosFunctionResult), module.WithValueSpec(a.sessionSpec), module.WithValueSpec(a.genericSpec)); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(a.ctx)

	go func() {
		defer cancel()

		if err = a.m.RegisterAndListen(ctx); err != nil {
			log.Error(err)
		}

	}()

	if err = a.entrypoint(); err != nil {
		return
	}

	<-ctx.Done()

	return
}
