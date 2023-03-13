// Copyright 2023 NJWS Inc.

package monitor

import (
	"net/url"
	"strings"

	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent"

	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/koron/go-ssdp"
)

func (a *Agent) monitor() (err error) {
	log.Info("ssdp monitor")

	m := &ssdp.Monitor{
		Alive: a.onAlive,
		Bye:   a.onBye,
	}

	return m.Start()
}

func (a *Agent) onAlive(m *ssdp.AliveMessage) {
	if !strings.Contains(m.Type, "redfish-rest") {
		return
	}

	if err := a.createOrUpdateAliveMessage(m); err != nil {
		log.Error(err)
		return
	}
}

func (a *Agent) onBye(m *ssdp.ByeMessage) {
	if !strings.Contains(m.Type, "redfish-rest") {
		return
	}

	log.Infof("Bye: From=%s Type=%s USN=%s", m.From.String(), m.Type, m.USN)
}

func (a *Agent) createOrUpdateAliveMessage(m *ssdp.AliveMessage) (err error) {
	redfishDevicesObject, err := a.getDocument(types.RedfishDevicesPath)
	if err != nil {
		return
	}

	description, err := device.GetDescription(m.Location)
	if err != nil {
		return
	}

	u, err := url.Parse(description.Device.PresentationURL)
	if err != nil {
		return
	}

	// FIXME fix if http
	u.Scheme = "https"

	redfishDevice := description.ToDevice(u)

	functionContext, err := agent.PrepareDiscoveryFunc(redfishDevicesObject.Id.String(), redfishDevice)
	if err != nil {
		return
	}

	return a.executor.ExecAsync(a.ctx, functionContext)
}
