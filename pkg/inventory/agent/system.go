// Copyright 2023 NJWS Inc.

package agent

import (
	"fmt"

	"git.fg-tech.ru/listware/cmdb/pkg/cmdb/documents"
	"git.fg-tech.ru/listware/go-core/pkg/client/system"
	"git.fg-tech.ru/listware/go-core/pkg/module"
	"git.fg-tech.ru/listware/proto/sdk/pbtypes"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/bootstrap"
	"github.com/stmcginnis/gofish/redfish"
)

// flink function
// systemsFunction (cmdbContext -> "%s.redfish-devices.root")
func (a *Agent) systemFunction(ctx module.Context) (err error) {
	client, err := a.getClientFromContext(ctx)
	if err != nil {
		return
	}

	m, err := a.getMessageFromContext(ctx)
	if err != nil {
		return
	}

	computerSystem, err := redfish.GetComputerSystem(client, m.Additional)
	if err != nil {
		return
	}
	return a.createOrUpdateSystem(ctx, computerSystem)
}

func (a *Agent) createOrUpdateSystem(ctx module.Context, computerSystem *redfish.ComputerSystem) (err error) {
	systemLink := fmt.Sprintf("system-%s", computerSystem.UUID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.*[?@._id == '%s'?].objects.root", systemLink, ctx.Self().Id)
	if err != nil {
		if functionContext, err = system.CreateChild(ctx.Self().Id, types.RedfishSystemID, systemLink, computerSystem); err != nil {
			return err
		}
	} else {
		if functionContext, err = system.UpdateObject(document.Id.String(), computerSystem); err != nil {
			return err
		}
	}

	functionContext.ReplyResult = ctx.GetReplyResult()

	m := make(map[string]string)
	ctx.Storage().Get(a.genericSpec, &m)

	m[functionContext.ReplyResult.GetKey()] = systemLink

	ctx.Storage().Set(a.genericSpec, m)

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}
	ctx.Send(msg)

	computerSystem.Bios()

	return // a.createOrUpdateSystemBIOS(ctx, redfishDevice, computerSystem)
}

func (a *Agent) onSystemFunctionResult(ctx module.Context, functionResult *pbtypes.FunctionResult) {
	if !functionResult.Complete {
		return
	}

	if err := a.inventoryBios(ctx, functionResult.GetReplyEgress().GetKey()); err != nil {
		ctx.AddError(err)
	}

}

func (a *Agent) inventoryBios(ctx module.Context, key string) (err error) {
	// m := make(map[string]string)
	// ctx.Storage().Get(a.genericSpec, &m)
	// systemLink := m[key]
	// delete(m, key)
	// ctx.Storage().Set(a.genericSpec, m)

	// document, err := a.getDocument("%s.*[?@._id == '%s'?].objects.root", systemLink, ctx.Self().Id)
	// if err != nil {
	// 	return
	// }

	// msg.Additional = link

	// data, err := msg.Marshal()
	// if err != nil {
	// 	return err
	// }

	// functionContext, err := prepareBiosFunc(document.Id.String(), data)
	// if err != nil {
	// 	return
	// }

	// msg, err := module.ToMessage(functionContext)
	// if err != nil {
	// 	return
	// }
	// ctx.Send(msg)
	return
}

func (a *Agent) createOrUpdateSystemBIOS(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	bios, err := computerSystem.Bios()
	if err != nil {
		return
	}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("bios.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-bios", "bios", bios)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), bios)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return
}

func (a *Agent) createOrUpdateSystemLed(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	led := &bootstrap.RedfishLed{Led: computerSystem.IndicatorLED}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("led.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-led", "led", led)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), led)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateSystemStatus(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemStatus(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	status := &bootstrap.RedfishStatus{Status: computerSystem.Status}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-status", "status", status)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), status)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateSystemBoot(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemBoot(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	boot := computerSystem.Boot

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("boot.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-boot", "boot", boot)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), boot)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateSystemBootOptions(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemBootOptions(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	parentNode, err := a.getDocument("boot.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	bootOptions, err := computerSystem.BootOptions()
	if err != nil {
		return
	}

	for _, bootOption := range bootOptions {
		if err = a.createOrUpdateBootOption(ctx, redfishDevice, parentNode, computerSystem, bootOption); err != nil {
			return
		}
	}
	return a.createOrUpdateSystemSecureBoot(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateBootOption(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, bootOption *redfish.BootOption) (err error) {
	bootOptionLink := fmt.Sprintf("option-%s", bootOption.ID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.boot.system-%s.service.%s.redfish-devices.root", bootOptionLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-boot-option", bootOptionLink, bootOption)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), bootOption)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return
}

func (a *Agent) createOrUpdateSystemSecureBoot(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	secureBoot, err := computerSystem.SecureBoot()
	if err != nil {
		return
	}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("secure-boot.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-secure-boot", "secure-boot", secureBoot)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), secureBoot)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdatePCIeDevices(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdatePCIeDevices(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	devices, err := computerSystem.PCIeDevices()
	if err != nil {
		return
	}

	for _, device := range devices {
		if err = a.createOrUpdatePCIeDevice(ctx, redfishDevice, parentNode, computerSystem, device); err != nil {
			return
		}
	}

	return a.createOrUpdateSystemPowerState(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdatePCIeDevice(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, device *redfish.PCIeDevice) (err error) {
	deviceLink := fmt.Sprintf("pcie-device-%s", device.ID)

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-pcie-device", deviceLink, device)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), device)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdatePCIeDeviceStatus(ctx, redfishDevice, computerSystem, deviceLink, device)
}

func (a *Agent) createOrUpdatePCIeDeviceStatus(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem, deviceLink string, device *redfish.PCIeDevice) (err error) {
	status := &bootstrap.RedfishStatus{Status: device.Status}

	parentNode, err := a.getDocument("%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-status", "status", status)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), status)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdatePCIeDeviceInterface(ctx, redfishDevice, computerSystem, deviceLink, device)
}

func (a *Agent) createOrUpdatePCIeDeviceInterface(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem, deviceLink string, device *redfish.PCIeDevice) (err error) {
	deviceInterface := &device.PCIeInterface

	parentNode, err := a.getDocument("%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("pcie-interface.%s.system-%s.service.%s.redfish-devices.root", deviceLink, computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-pcie-interface", "pcie-interface", deviceInterface)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), deviceInterface)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return
}

func (a *Agent) createOrUpdateSystemPowerState(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	powerState := &bootstrap.RedfishPowerState{PowerState: computerSystem.PowerState}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("power-state.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-power-state", "power-state", powerState)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), powerState)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateSystemPowerRestorePolicy(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateSystemPowerRestorePolicy(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	powerRestorePolicy := &bootstrap.RedfishPowerRestorePolicy{PowerRestorePolicy: computerSystem.PowerRestorePolicy}

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("power-restore-policy.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-power-restore-policy", "power-restore-policy", powerRestorePolicy)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), powerRestorePolicy)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateProcessorSummary(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateProcessorSummary(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	processorSummary := computerSystem.ProcessorSummary

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("processor-summary.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-processor-summary", "processor-summary", processorSummary)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), processorSummary)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateProcessorSummaryStatus(ctx, redfishDevice, computerSystem, &processorSummary)
}

func (a *Agent) createOrUpdateProcessorSummaryStatus(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem, processorSummary *redfish.ProcessorSummary) (err error) {
	status := &bootstrap.RedfishStatus{Status: processorSummary.Status}

	parentNode, err := a.getDocument("processor-summary.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.processor-summary.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-status", "status", status)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), status)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateMemorySummary(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateMemorySummary(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	memorySummary := computerSystem.MemorySummary

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("memory-summary.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-memory-summary", "memory-summary", memorySummary)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), memorySummary)
		if err != nil {
			return err
		}
	}

	if err = a.executor.ExecSync(ctx, functionContext); err != nil {
		return
	}

	return a.createOrUpdateMemorySummaryStatus(ctx, redfishDevice, computerSystem, &memorySummary)
}

func (a *Agent) createOrUpdateMemorySummaryStatus(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem, memorySummary *redfish.MemorySummary) (err error) {
	status := &bootstrap.RedfishStatus{Status: memorySummary.Status}

	parentNode, err := a.getDocument("memory-summary.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("status.memory-summary.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-status", "status", status)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), status)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	return a.createOrUpdateHostWatchdogTimer(ctx, redfishDevice, computerSystem)
}

func (a *Agent) createOrUpdateHostWatchdogTimer(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	hostWatchdogTimer := computerSystem.HostWatchdogTimer

	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	var functionContext *pbtypes.FunctionContext
	document, err := a.getDocument("host-watchdog-timer.system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		functionContext, err = system.CreateChild(parentNode.Id.String(), "types/redfish-host-watchdog-timer", "host-watchdog-timer", hostWatchdogTimer)
		if err != nil {
			return err
		}
	} else {
		functionContext, err = system.UpdateObject(document.Id.String(), hostWatchdogTimer)
		if err != nil {
			return err
		}
	}

	msg, err := module.ToMessage(functionContext)
	if err != nil {
		return
	}

	ctx.Send(msg)

	// TODO: return createMemories, createProcessors, etc.
	return
}

// TODO: currently unavailable structures, pending to implement
func (a *Agent) createOrUpdateMemories(ctx module.Context, redfishDevice device.RedfishDevice, computerSystem *redfish.ComputerSystem) (err error) {
	parentNode, err := a.getDocument("system-%s.service.%s.redfish-devices.root", computerSystem.UUID, redfishDevice.UUID())
	if err != nil {
		return
	}

	memories, err := computerSystem.Memory()
	if err != nil {
		return
	}

	for _, memory := range memories {
		if err = a.createOrUpdateMemory(ctx, redfishDevice, parentNode, computerSystem, memory); err != nil {
			return
		}
	}
	return
}

func (a *Agent) createOrUpdateMemory(ctx module.Context, redfishDevice device.RedfishDevice, parentNode *documents.Node, computerSystem *redfish.ComputerSystem, memory *redfish.Memory) (err error) {
	fmt.Println(memory.DeviceLocator, memory.MemoryLocation)
	return
}
