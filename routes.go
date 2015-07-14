package main

import (
	"net/http"

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

type Routes []Route

var routes = Routes{
	Route{
		"group_availability",
		"GET",
		"group_availability/",
		Respond(endpointGroupAvailability.List),
	},
	Route{
		"group_groups_availability",
		"GET",
		"group_groups_availability",
		Respond(groupGroupsAvailability.List),
	},
	Route{
		"endpoint_group_availability",
		"POST",
		"endpoint_group_availability",
		Respond(endpointGroupAvailability.List),
	},
	Route{
		"service_flavor_availability",
		"GET",
		"service_flavor_availability",
		Respond(serviceFlavorAvailability.List),
	},
	Route{
		"AP List",
		"GET",
		"AP",
		Respond(availabilityProfiles.List),
	},
	Route{
		"AP Create",
		"POST",
		"AP",
		Respond(availabilityProfiles.Create),
	},
	Route{
		"AP update",
		"PUT",
		"AP/{id}",
		Respond(availabilityProfiles.Update),
	},
	Route{
		"AP delete",
		"DELETE",
		"AP/{id}",
		Respond(availabilityProfiles.Delete),
	},
	Route{
		"PLACEHOLDER",
		"GET",
		"service_flavor_availability",
		Respond(serviceFlavorAvailability.List),
	},
	Route{
		"tenant create",
		"GET",
		"tenants",
		Respond(tenants.Create),
	},
	Route{
		"tenant update",
		"PUT",
		"tenants/{name}",
		Respond(tenants.Update),
	},
	Route{
		"tenant delete",
		"DELETE",
		"tenants/{name}",
		Respond(tenants.Delete),
	},
	Route{
		"tenant list",
		"GET",
		"/tenants",
		Respond(tenants.List),
	},
	Route{
		"tenant list one",
		"GET",
		"/tenants/{name}",
		Respond(tenants.ListOne),
	},

	//jobs
	Route{
		"jobs create",
		"POST",
		"jobs",
		Respond(jobs.Create),
	},
	Route{
		"job update",
		"PUT",
		"jobs/{name}",
		Respond(jobs.Update),
	},
	Route{
		"job delete",
		"DELETE",
		"jobs/{name}",
		Respond(jobs.Delete),
	},
	Route{
		"job list",
		"GET",
		"jobs",
		Respond(jobs.List),
	},
	Route{
		"job list one",
		"GET",
		"jobs/{name}",
		Respond(jobs.ListOne),
	},

	//Poem Profiles compatibility
	Route{
		"List poems",
		"GET",
		"poems",
		Respond(metricProfiles.ListPoems),
	},

	//Metric Profiles
	Route{
		"list metric profile",
		"GET",
		"metric_profiles",
		Respond(metricProfiles.List),
	},
	Route{
		"metric profile create",
		"POST",
		"metric_profiles",
		Respond(metricProfiles.Create),
	},
	Route{
		"metric profile delete",
		"DELETE",
		"metric_profiles/{id}",
		Respond(metricProfiles.Delete),
	},
	Route{
		"metric profile update",
		"PUT",
		"metric_profiles/{id}",
		Respond(metricProfiles.Update),
	},

	//Recalculations
	Route{
		"recomputation create",
		"POST",
		"recomputations",
		Respond(recomputations.Create),
	},
	Route{
		"recomputation list",
		"GET",
		"recomputations",
		Respond(recomputations.List),
	},

	Route{
		"factors list",
		"GET",
		"factors",
		Respond(factors.List),
	},

	//Status
	Route{
		"status detail list",
		"GET",
		"status/metrics/timeline/{group}",
		Respond(statusDetail.List),
	},

	//Status Raw Msg
	Route{
		"status message list",
		"GET",
		"status/metrics/msg/{hostname}/{service}/{metric}",
		Respond(statusMsg.List),
	},

	//Status Endpoints
	Route{
		"status endpoint list",
		"GET",
		"status/endpoints/timeline/{hostname}/{service_type}",
		Respond(statusEndpoints.List),
	},

	//Status Services
	Route{
		"status service list",
		"GET",
		"status/services/timeline/{group}",
		Respond(statusServices.List),
	},

	//Status Sites
	Route{
		"status endpoint group list",
		"GET",
		"status/sites/timeline/{group}",
		Respond(statusEndpointGroups.List),
	},
}
