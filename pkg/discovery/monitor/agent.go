// Copyright 2023 NJWS Inc.

package monitor

import (
	"context"

	"git.fg-tech.ru/listware/go-core/pkg/executor"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.New()
)

type Agent struct {
	ctx    context.Context
	cancel context.CancelFunc

	executor executor.Executor
}

// Run agent
func Run(ctx context.Context) (err error) {
	a := &Agent{}
	a.ctx, a.cancel = context.WithCancel(ctx)

	if a.executor, err = executor.New(); err != nil {
		return
	}

	return a.run()
}

func (a *Agent) run() (err error) {
	defer a.executor.Close()

	log.Infof("run monitor agent")

	a.osSignalCtrl()

	ctx, cancel := context.WithCancel(a.ctx)
	defer cancel()

	if err = a.monitor(); err != nil {
		return
	}

	<-ctx.Done()

	return
}
