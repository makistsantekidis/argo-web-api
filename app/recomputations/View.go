/*
 * Copyright (c) 2015 GRNET S.A.
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
 * or implied, of GRNET S.A.
 *
 */

package recomputations

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/ARGOeu/argo-web-api/respond"
)

func createListView(results interface{}, format string) ([]byte, error) {

	docRoot := &respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: "Success",
			Code:    "200",
		},
	}

	docRoot.Data = results
	if format == "application/xml" {
		output, err := xml.MarshalIndent(docRoot, "", " ")
		return output, err
	}

	output, err := json.MarshalIndent(docRoot, "", " ")
	return output, err
}

func createSubmitView(inserted MongoInterface, format string, r *http.Request) ([]byte, error) {
	docRoot := &respond.ResponseMessage{
		Status: respond.StatusResponse{
			Message: "Recomputations successfully created",
			Code:    "201",
		},
		Data: SelfReference{
			ID:    inserted.ID,
			Links: Links{Self: "https://" + r.Host + r.URL.Path + "/" + inserted.ID},
		},
	}

	// Message{
	// 	Message: "Recomputations successfully submitted",
	// 	Status:  "202",
	// }

	output, err := json.MarshalIndent(docRoot, "", " ")
	return output, err
}

func messageXML(answer string) ([]byte, error) {
	docRoot := &Message{}
	docRoot.Message = answer
	output, err := xml.MarshalIndent(docRoot, " ", "  ")
	return output, err
}
