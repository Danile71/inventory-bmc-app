// Copyright 2023 NJWS Inc.

package monitor

import (
	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"github.com/foliagecp/inventory-bmc-app/pkg/utils"
)

func (a *Agent) getDocument(format string, args ...any) (document *documents.Node, err error) {
	return utils.GetDocument(a.ctx, format, args...)
}
