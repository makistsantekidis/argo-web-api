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
	"strings"

	"labix.org/v2/mgo/bson"

	"github.com/argoeu/argo-web-api/utils/authentication"
	"github.com/argoeu/argo-web-api/utils/config"
	"github.com/argoeu/argo-web-api/utils/mongo"

	"github.com/argoeu/argo-web-api/respond"
	"github.com/gorilla/mux"
)

// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter
func HandleSubrouter(s *mux.Router, confhandler *respond.ConfHandler) {
	// Route for request /api/v2/results/
	// s.Path("/").Handler(routing.Respond())

	// Route for request "/api/v2/results/{report_name}"
	reportSubrouter := s.PathPrefix("/{report_name}").Subrouter()
	// TODO: list reports with the name {reportn_name}
	// reportSubrouter.Name("reports").Handler(confhandler.Respond(ListReports))

	// Route for request "api/v2/results/{report_name}/{group_type}"
	groupTypeSubrouter := reportSubrouter.PathPrefix("/{group_type}").Subrouter()
	// TODO: list groups with type {group_type}
	// groupTypeSubrouter.
	// 	Path("/").
	// 	Name("Group_Type").
	// 	Handler(respond.Respond(...))
	// TODO: list endpointgroups with type {lgroup_type}
	// groupTypeSubrouter.
	// Path("/{group_name}/{lgroup_name}").
	// Name("EndpointGroup_Type").
	// Handler(respond.Respond(...))

	// Route for request "api/v2/results/{report_name}/{group_type}/{group_name}"
	// matches only endpoint groups
	groupSubrouter := groupTypeSubrouter.PathPrefix("/{group_name}").Subrouter()
	groupSubrouter.
		Methods("GET").
		Name("Group Name").
		MatcherFunc(matchEndpointGroup(confhandler.Config)).
		Handler(confhandler.Respond(ListEndpointGroupResults))

	// groupSubrouter.
	// 	Methods("GET").
	// 	Name("Group Name").
	// 	Handler(confhandler.Respond(ListSuperGroupResults, "group name"))

	groupSubrouter.Methods("GET").
		Path("/{lgroup_type}/{lgroup_name}").
		Name("Group/LGroup Names").
		Handler(confhandler.Respond(ListEndpointGroupResults))

	// Route for request "api/v2/results/{report_name}/{group_type}/{group_name}/services"
	// matches only endpoint groups
	serviceSubrouter := groupSubrouter.PathPrefix("/services").Subrouter()
	
	serviceSubrouter.
		Methods("GET").
		Name("Services").
		Handler(confhandler.Respond(ListServiceFlavorResults))

	serviceSubrouter.
		Path("/{service_flavor}").
		Methods("GET").
		Name("Services/Service Flavor").
		Handler(confhandler.Respond(ListServiceFlavorResults))

}

func matchEndpointGroup(cfg config.Config) mux.MatcherFunc {
	return func(r *http.Request, routematch *mux.RouteMatch) bool {
		// vars := mux.Vars(r)
		// name := vars["group_name"]
		// or
		// name := routematch.Vars
		// But for some reason none of these work so this is a workaround:
		// TODO: create reproduction code and open issue on gorilla/mux issue tracker

		name := strings.Split(strings.Split(strings.Split(r.URL.String(), "results/")[1], "/")[2], "?")[0]
		tenantcfg, err := authentication.AuthenticateTenant(r.Header, cfg)
		if err != nil {
			return false
		}
		session, err := mongo.OpenSession(tenantcfg)
		if err != nil {
			return false
		}
		result := bson.M{}
		err = mongo.FindOne(session, tenantcfg.Db, "endpoint_group_ar", bson.M{"name": name}, result)
		if err != nil {
			return false
		}
		return true
	}
}
