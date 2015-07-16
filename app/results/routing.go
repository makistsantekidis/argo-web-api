package results

import (

    "github.com/argoeu/argo-web-api/utils/logging"

	"github.com/gorilla/mux"

)

func HandleSubrouter(s *mux.Router){
    reportSubrouter := s.PathPrefix("/{report_name}").Subrouter()



}
