// Copyright 2023 NJWS Inc.

package agent

import (
	"context"
	"fmt"

	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
)

func createOrUpdateFunctionLink(ctx context.Context, fromQuery, toQuery, name string) (functionContext *pbtypes.FunctionContext, err error) {
	route := &pbtypes.FunctionRoute{
		Url: "http://inventory-bmc:31001/statefun",
	}

	query := fmt.Sprintf("%s.%s", name, fromQuery)

	fmt.Println(query)

	if linkDocument, err := getDocument(ctx, query); err == nil {
		return system.UpdateAdvancedLink(linkDocument.LinkId.String(), route)
	}

	parent, err := getDocument(ctx, fromQuery)
	if err != nil {
		return
	}
	fmt.Println("creates", fromQuery)

	child, err := getDocument(ctx, toQuery)
	if err != nil {
		return
	}

	return system.CreateLink(parent.Id.String(), child.Id.String(), name, "function", route)
}

func (a *Agent) createOrUpdateFunctionLink(fromQuery, toQuery, name string) (err error) {
	functionContext, err := createOrUpdateFunctionLink(a.ctx, fromQuery, toQuery, name)
	if err != nil {
		return
	}

	return a.executor.ExecSync(a.ctx, functionContext)
}

/*
func createOrUpdateFunctionLink(ctx context.Context, id documents.DocumentID, functionPath string) (functionContext *pbtypes.FunctionContext, err error) {
	route := &pbtypes.FunctionRoute{
		Url:             "http://inventory-bmc:31001/statefun",
		ExecuteOnCreate: true,
		ExecuteOnUpdate: true,
	}

	query := fmt.Sprintf("%s.%s", id.Key(), functionPath)

	if linkDocument, err := getDocument(ctx, query); err == nil {
		return system.UpdateAdvancedLink(linkDocument.LinkId.String(), route)
	}

	function, err := getDocument(ctx, functionPath)
	if err != nil {
		return
	}

	return system.CreateLink(function.Id.String(), id.String(), id.Key(), function.Type, route)
}

func CreateOrUpdateInventoryLink(ctx context.Context, id documents.DocumentID) (functionContext *pbtypes.FunctionContext, err error) {
	return createOrUpdateFunctionLink(ctx, id, types.InventoryFunctionPath)
}

func CreateOrUpdateServiceLink(ctx context.Context, id documents.DocumentID) (functionContext *pbtypes.FunctionContext, err error) {
	return createOrUpdateFunctionLink(ctx, id, types.ServiceFunctionPath)
}

func CreateOrUpdateSystemsLink(ctx context.Context, id documents.DocumentID) (functionContext *pbtypes.FunctionContext, err error) {
	return createOrUpdateFunctionLink(ctx, id, types.SystemsFunctionPath)
}

func CreateOrUpdateSystemLink(ctx context.Context, id documents.DocumentID) (functionContext *pbtypes.FunctionContext, err error) {
	return createOrUpdateFunctionLink(ctx, id, types.SystemFunctionPath)
}
*/
