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

package results

import (
	"net/http"

	"github.com/argoeu/argo-web-api/utils/config"

	"github.com/argoeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, apphandler *respond.AppHandler) {
	// Route for request /api/v2/results/
	// s.Path("/").Handler(routing.Respond())

	// Route for request "/api/v2/results/{report_name}"
	reportSubrouter := s.PathPrefix("/{report_name}").Subrouter()
	// TODO: list reports with the name {reportn_name}
	// reportSubrouter.Path("/").Name("reports").Handler(respond.Respond(ListReports, "reports", cfg))

	// Route for request "api/v2/results/{report_name}/{group_type}"
	groupTypeSubrouter := reportSubrouter.PathPrefix("/{group_type}").Subrouter()
	// TODO: list groups with type {group_type}
	// groupTypeSubrouter.Path("/").Name("Group_Type").Handler(respond.Respond(List))

	// Route for request "api/v2/results/{report_name}/{group_type}/{group_name}"
	groupSubrouter := groupTypeSubrouter.PathPrefix("/{group_name}").Subrouter()
	groupSubrouter.Path("/").
		Name("Group Name").
		MatcherFunc(MatchEndpointGroup(apphandler.Cfg)).
		Handler(apphandler.Respond(ListEndpointGroupResults, "group name"))

}

func MatchEndpointGroup(cfg config.Config) mux.MatcherFunc {
	return func(r *http.Request, routematch *mux.RouteMatch) bool {
		// vars := mux.Vars(r)
		//TODO: figure out if vars['group_type'] is group or endpointgroup
		return true
	}
}
