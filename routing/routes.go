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
    "net/http"
    "github.com/gorilla/mux"

    "github.com/argoeu/argo-web-api/app/results"
    "github.com/argoeu/argo-web-api/app/availabilityProfiles"
    "github.com/argoeu/argo-web-api/app/endpointGroupAvailability"
    "github.com/argoeu/argo-web-api/app/factors"
    "github.com/argoeu/argo-web-api/app/groupGroupsAvailability"
    "github.com/argoeu/argo-web-api/app/jobs"
    "github.com/argoeu/argo-web-api/app/metricProfiles"
    "github.com/argoeu/argo-web-api/app/recomputations"
    "github.com/argoeu/argo-web-api/app/serviceFlavorAvailability"
    "github.com/argoeu/argo-web-api/app/statusDetail"
    "github.com/argoeu/argo-web-api/app/statusEndpointGroups"
    "github.com/argoeu/argo-web-api/app/statusEndpoints"
    "github.com/argoeu/argo-web-api/app/statusMsg"
    "github.com/argoeu/argo-web-api/app/statusServices"
    "github.com/argoeu/argo-web-api/app/tenants"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type SubRoute struct {
    Name string
    Pattern string
    SubrouterHandler func(*mux.Router)
}

type Routes []Route
type SubRoutes []SubRoute

var subroutes = SubRoutes{
    {"results", "/results", results.HandleSubrouter},
    
}

var routes = Routes{

    //-----------------------------------New requests, need to fill in the Respond(...) functions---------------------------------------------
    // {"groups of specific type",     "GET",      "/results/{report}/{group_type}",                                                   Respond()},
    // {"Group of groups",             "GET",      "/results/{report}/{group_type}/{group}",                                           Respond()},
    // {"Endpoint group in group",     "GET",      "/results/{report}/{group_type}/{group}/{lgroup_type}/{lgroup}",                    Respond()},
    // {"All services of group",       "GET",      "/results/{report}/{group_type}/{group}/{lgroup_type}/{lgroup}/services",           Respond()},
    // {"Specific service of group",   "GET",      "/results/{report}/{group_type}/{group}/{lgroup_type}/{lgroup}/services/{service}", Respond()},
    // {"All services of group",       "GET",      "/results/{report}/{lgroup_type}/{lgroup}/services",                                Respond()},
    // {"Specific service of group",   "GET",      "/results/{report}/{lgroup_type}/{lgroup}/services/{service}",                      Respond()},


    //-----------------------------------Old requests for here on down -------------------------------------------------
    {"group_availability",          "GET",      "/group_availability/",         Respond(endpointGroupAvailability.List)},
    {"group_groups_availability",   "GET",      "/group_groups_availability",   Respond(groupGroupsAvailability.List)},
    {"endpoint_group_availability", "POST",     "/endpoint_group_availability", Respond(endpointGroupAvailability.List)},
    {"service_flavor_availability", "GET",      "/service_flavor_availability", Respond(serviceFlavorAvailability.List)},
    {"AP List",                     "GET",      "/AP",                          Respond(availabilityProfiles.List)},
    {"AP Create",                   "POST",     "/AP",                          Respond(availabilityProfiles.Create)},
    {"AP update",                   "PUT",      "/AP/{id}",                     Respond(availabilityProfiles.Update)},
    {"AP delete",                   "DELETE",   "/AP/{id}",                     Respond(availabilityProfiles.Delete)},
    {"PLACEHOLDER",                 "GET",      "/service_flavor_availability", Respond(serviceFlavorAvailability.List)},
    {"tenant create",               "GET",      "/tenants",                     Respond(tenants.Create)},
    {"tenant update",               "PUT",      "/tenants/{name}",              Respond(tenants.Update)},
    {"tenant delete",               "DELETE",   "/tenants/{name}",              Respond(tenants.Delete)},
    {"tenant list",                 "GET",      "/tenants",                     Respond(tenants.List)},
    {"tenant list one",             "GET",      "/tenants/{name}",              Respond(tenants.ListOne)},

    //jobs
    {"jobs create",                 "POST",     "/jobs",                        Respond(jobs.Create)},
    {"job update",                  "PUT",      "/jobs/{name}",                 Respond(jobs.Update)},
    {"job delete",                  "DELETE",   "/jobs/{name}",                 Respond(jobs.Delete)},
    {"job list",                    "GET",      "/jobs",                        Respond(jobs.List)},
    {"job list one",                "GET",      "/jobs/{name}",                 Respond(jobs.ListOne)},

    //Poem Profiles compatibility
    {"List poems", "GET", "/poems", Respond(metricProfiles.ListPoems)},

    //Metric Profiles
    {"list metric profile",     "GET",      "/metric_profiles",                 Respond(metricProfiles.List)},
    {"metric profile create",   "POST",     "/metric_profiles",                 Respond(metricProfiles.Create)},
    {"metric profile delete",   "DELETE",   "/metric_profiles/{id}",            Respond(metricProfiles.Delete)},
    {"metric profile update",   "PUT",      "/metric_profiles/{id}",            Respond(metricProfiles.Update)},

    //Recalculations
    {"recomputation create",    "POST", "/recomputations", Respond(recomputations.Create)},
    {"recomputation list",      "GET",  "/recomputations", Respond(recomputations.List)},

    {"factors list", "GET", "factors", Respond(factors.List)},

    //Status
    {"status detail list", "GET", "status/metrics/timeline/{group}", Respond(statusDetail.List)},

    //Status Raw Msg
    {"status message list", "GET", "status/metrics/msg/{hostname}/{service}/{metric}", Respond(statusMsg.List)},

    //Status Endpoints
    {"status endpoint list", "GET", "status/endpoints/timeline/{hostname}/{service_type}", Respond(statusEndpoints.List)},

    //Status Services
    {"status service list", "GET", "status/services/timeline/{group}", Respond(statusServices.List)},

    //Status Sites
    {"status endpoint group list", "GET", "status/sites/timeline/{group}", Respond(statusEndpointGroups.List)},
}
