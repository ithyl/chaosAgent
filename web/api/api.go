/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"chaosAgent/transport"
	chaosWeb "chaosAgent/web"
	"chaosAgent/web/handler"
	"chaosAgent/web/server"
)

type API struct {
	chaosWeb.APiServer
	//ready func(http.HandlerFunc) http.HandlerFunc

}

// community just use http
func NewAPI() *API {

	return &API{
		server.NewHttpServer(),
	}
}

func (api *API) Register(transportClient *transport.TransportClient) error {

	chaosbladeHandler := NewServerRequestHandler(handler.NewChaosbladeHandler(transportClient))
	if err := api.RegisterHandler("chaosblade", chaosbladeHandler); err != nil {
		return err
	}

	pingHandler := NewServerRequestHandler(handler.NewPingHandler())
	if err := api.RegisterHandler("ping", pingHandler); err != nil {
		return err
	}

	uninstallHandler := NewServerRequestHandler(handler.NewUninstallInstallHandler(transportClient))
	if err := api.RegisterHandler("uninstall", uninstallHandler); err != nil {
		return err
	}

	return nil
}
