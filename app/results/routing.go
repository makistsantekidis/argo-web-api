package results

import (

    "github.com/argoeu/argo-web-api/utils/logging"
    "github.com/argoeu/argo-web-api/routing"


	"github.com/gorilla/mux"

)
// HandleSubrouter uses the subrouter for a specific calls and creates a tree of sorts
// handling each route with a different subrouter 
func HandleSubrouter(s *mux.Router){

    // Route for request "/api/v2/results/{report_name}"
    reportSubrouter := s.PathPrefix("/{report_name}").Subrouter()
    reportSubrouter.Path("/").Name("reports").Handler(routing.Respond(List))

    // Route for request "api/v2/results/{report_name}/{group_type}"
    groupSubrouter := reportSubrouter.PathPrefix("/{group_type}").Subrouter()
    groupSubrouter.Path("/").Name("Group_Type").Handler(routing(Respond(List)))

}
