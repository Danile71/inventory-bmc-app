// Copyright 2023 NJWS Inc.

package types

const (
	App       = "inventory-bmc"
	Namespace = "proxy.foliage"

	// FunctionType will be as default qdsl to function
	InventoryFunctionType = InventoryFunctionPath
	ServiceFunctionType   = ServiceFunctionPath
	SystemsFunctionType   = SystemsFunctionPath
	SystemFunctionType    = SystemFunctionPath
	BiosFunctionType      = BiosFunctionPath

	Description = "inventory redfish function"
)

const (
	FunctionContainerID = "types/function-container"
	FunctionID          = "types/function"
	RedfishServiceID    = "types/redfish-service"
	RedfishSystemID     = "types/redfish-system"
	RootID              = "system/root"
)

const (
	RedfishDeviceKey      = "redfish-device"
	FunctionContainerLink = "redfish"
	InventoryFunctionLink = "inventory"
	ServiceFunctionLink   = "service"
	SystemsFunctionLink   = "systems"
	SystemFunctionLink    = "system"
	BiosFunctionLink      = "bios"

	RedfishServiceLink = "service"
)

const (
	FunctionsPath         = "functions.root"
	FunctionContainerPath = "redfish.functions.root"
	InventoryFunctionPath = "inventory.redfish.functions.root"
	ServiceFunctionPath   = "service.inventory.redfish.functions.root"
	SystemsFunctionPath   = "systems.inventory.redfish.functions.root"
	SystemFunctionPath    = "system.inventory.redfish.functions.root"
	BiosFunctionPath      = "bios.inventory.redfish.functions.root"
)
