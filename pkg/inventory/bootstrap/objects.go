// Copyright 2023 NJWS Inc.

package bootstrap

import (
	"context"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/qdsl"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/proto/sdk/pbcmdb"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
)

var (
	registerObjects = []*pbcmdb.RegisterObjectMessage{}
)

type RedfishFunctionContainer struct{}

func createObjects(ctx context.Context) (err error) {
	if err = createRedFishMountpointObject(ctx); err != nil {
		return
	}

	if err = createInventoryFunctionObject(ctx); err != nil {
		return
	}

	if err = createServiceFunctionObject(ctx); err != nil {
		return
	}

	if err = createSystemsFunctionObject(ctx); err != nil {
		return
	}

	if err = createSystemFunctionObject(ctx); err != nil {
		return
	}

	if err = createBiosFunctionObject(ctx); err != nil {
		return
	}

	return
}

func registerObject(ctx context.Context, objectPath string, message *pbcmdb.RegisterObjectMessage) (err error) {
	// check if object exists
	elements, err := qdsl.Qdsl(ctx, objectPath)
	if err != nil {
		return
	}

	// TODO already exists
	if len(elements) > 0 {
		return
	}

	registerObjects = append(registerObjects, message)
	return
}

func createRedFishMountpointObject(ctx context.Context) (err error) {
	message, err := system.RegisterObject(types.FunctionsPath, types.FunctionContainerID, types.FunctionContainerLink, RedfishFunctionContainer{}, false, true)
	if err != nil {
		return
	}
	return registerObject(ctx, types.FunctionContainerPath, message)
}

func createInventoryFunctionObject(ctx context.Context) (err error) {
	function := pbtypes.Function{
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.InventoryFunctionType,
		},
		Description: types.Description,
		Grounded:    false,
	}

	message, err := system.RegisterObject(types.FunctionContainerPath, types.FunctionID, types.InventoryFunctionLink, function, false, true)
	if err != nil {
		return
	}
	return registerObject(ctx, types.InventoryFunctionPath, message)
}

func createServiceFunctionObject(ctx context.Context) (err error) {
	function := pbtypes.Function{
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.ServiceFunctionType,
		},
		Description: types.Description,
		Grounded:    false,
	}

	message, err := system.RegisterObject(types.InventoryFunctionPath, types.FunctionID, types.ServiceFunctionLink, function, true, true)
	if err != nil {
		return
	}
	return registerObject(ctx, types.ServiceFunctionPath, message)
}

func createSystemsFunctionObject(ctx context.Context) (err error) {
	function := pbtypes.Function{
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.SystemsFunctionType,
		},
		Description: types.Description,
		Grounded:    false,
	}

	message, err := system.RegisterObject(types.InventoryFunctionPath, types.FunctionID, types.SystemsFunctionLink, function, true, true)
	if err != nil {
		return
	}
	return registerObject(ctx, types.SystemsFunctionPath, message)

}

func createSystemFunctionObject(ctx context.Context) (err error) {
	function := pbtypes.Function{
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.SystemFunctionType,
		},
		Description: types.Description,
		Grounded:    false,
	}

	message, err := system.RegisterObject(types.InventoryFunctionPath, types.FunctionID, types.SystemFunctionLink, function, true, true)
	if err != nil {
		return
	}
	return registerObject(ctx, types.SystemFunctionPath, message)
}

func createBiosFunctionObject(ctx context.Context) (err error) {
	function := pbtypes.Function{
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.BiosFunctionType,
		},
		Description: types.Description,
		Grounded:    false,
	}

	message, err := system.RegisterObject(types.InventoryFunctionPath, types.FunctionID, types.BiosFunctionLink, function, true, true)
	if err != nil {
		return
	}
	return registerObject(ctx, types.BiosFunctionPath, message)
}
