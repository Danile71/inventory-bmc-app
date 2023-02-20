// Copyright 2023 NJWS Inc.

package agent

import (
	"encoding/json"

	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types/message"
)

// type Request struct {
// 	Query string `json:"query"`
// 	Name  string `json:"name"`
// }

// func prepareFunc(id string, r Request) (fc *pbtypes.FunctionContext, err error) {
// 	ft := &pbtypes.FunctionType{
// 		Namespace: types.Namespace,
// 		Type:      types.FunctionPath,
// 	}

// 	fc = &pbtypes.FunctionContext{
// 		Id:           id,
// 		FunctionType: ft,
// 	}
// 	fc.Value, err = json.Marshal(r)
// 	return
// }

// // genFunction generate function call with object uuid and qdsl
// func genFunction(id, query string) (*pbtypes.FunctionContext, error) {
// 	r := Request{Query: query}
// 	return prepareFunc(id, r)
// }

func PrepareInventoryFunc(id string) (fc *pbtypes.FunctionContext, err error) {
	fc = &pbtypes.FunctionContext{
		Id: id,
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.InventoryFunctionPath,
		},
	}
	return
}

func prepareServiceFunc(id string) (fc *pbtypes.FunctionContext, err error) {
	fc = &pbtypes.FunctionContext{
		Id: id,
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.ServiceFunctionPath,
		},
	}
	return
}

func prepareSystemsFunc(id string, m *message.Message) (fc *pbtypes.FunctionContext, err error) {
	fc = &pbtypes.FunctionContext{
		Id: id,
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.SystemsFunctionPath,
		},
	}
	if fc.Value, err = json.Marshal(m); err != nil {
		return
	}
	return
}

func prepareSystemFunc(id string, m *message.Message) (fc *pbtypes.FunctionContext, err error) {
	fc = &pbtypes.FunctionContext{
		Id: id,
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.SystemFunctionPath,
		},
	}
	if fc.Value, err = json.Marshal(m); err != nil {
		return
	}
	return
}

func prepareBiosFunc(id string, m *message.Message) (fc *pbtypes.FunctionContext, err error) {
	fc = &pbtypes.FunctionContext{
		Id: id,
		FunctionType: &pbtypes.FunctionType{
			Namespace: types.Namespace,
			Type:      types.BiosFunctionPath,
		},
	}
	if fc.Value, err = json.Marshal(m); err != nil {
		return
	}
	return
}
