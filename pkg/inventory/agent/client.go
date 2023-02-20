// Copyright 2023 NJWS Inc.

package agent

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"git.fg-tech.ru/listware/go-core/pkg/module"
	"github.com/foliagecp/inventory-bmc-app/pkg/discovery/agent/types/redfish/device"
	"github.com/foliagecp/inventory-bmc-app/pkg/inventory/agent/types/message"
	"github.com/stmcginnis/gofish"
)

func (a *Agent) getClientFromCmdb(ctx module.Context) (c *gofish.APIClient, err error) {
	var redfishDevice device.RedfishDevice
	if err = json.Unmarshal(ctx.CmdbContext(), &redfishDevice); err != nil {
		return
	}

	u, err := url.Parse(redfishDevice.Api)
	if err != nil {
		return
	}

	config := gofish.ClientConfig{
		Endpoint: fmt.Sprintf("%s://%s", u.Scheme, u.Host),
		Username: redfishDevice.Login,
		Password: redfishDevice.Password,

		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}},
	}

	if c, err = gofish.ConnectContext(ctx, config); err != nil {
		return
	}

	// save to state

	m := &message.Message{
		Endpoint: fmt.Sprintf("%s://%s", u.Scheme, u.Host),
	}

	if m.Session, err = c.GetSession(); err != nil {
		return
	}

	ctx.Storage().Set(a.sessionSpec, m)
	return
}

func (a *Agent) getMessageFromContext(ctx module.Context) (m *message.Message, err error) {
	return message.Parse(ctx.Message())
}

func (a *Agent) getClientFromContext(ctx module.Context) (c *gofish.APIClient, err error) {
	m, err := a.getMessageFromContext(ctx)
	if err != nil {
		return
	}
	ctx.Storage().Set(a.sessionSpec, m)

	config := gofish.ClientConfig{
		Endpoint: m.Endpoint,
		Session:  m.Session,

		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}},
	}

	return gofish.ConnectContext(ctx, config)
}
