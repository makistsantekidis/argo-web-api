/*
 * Copyright (c) 2014 GRNET S.A., SRCE, IN2P3 CNRS Computing Centre
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the
 * License. You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an "AS
 * IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 *
 * The views and conclusions contained in the software and
 * documentation are those of the authors and should not be
 * interpreted as representing official policies, either expressed
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package routing

import (
	"github.com/argoeu/argo-web-api/respond"
	"github.com/argoeu/argo-web-api/utils/config"

	"github.com/gorilla/mux"
)

func NewRouter(cfg config.Config) *mux.Router {

	apphandler := respond.appHandler{cfg}
	router := mux.NewRouter() //.StrictSlash(true)
	for _, route := range routes {
		// var handler http.Handler

		handler := route.HandlerFunc
		handler = respond.Respond(handler, route.Name, cfg)

		router.
			PathPrefix("/api/v1").
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	router.Walk(apphandler.walker)

	for _, subroute := range subroutes {
		subrouter := router.
			PathPrefix("/api/v2").
			PathPrefix(subroute.Pattern).
			Subrouter()
		subroute.SubrouterHandler(subrouter, apphandler)
	}

	return router
}

func (apphandler *appHandler) walker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	route.Handler(apphandler.ServeHTTP(route.GetHandler()))
}
