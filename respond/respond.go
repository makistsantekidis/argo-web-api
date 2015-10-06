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
 * or implied, of either GRNET S.A., SRCE or IN2P3 CNRS Computing
 * Centre
 *
 * The work represented by this source file is partially funded by
 * the EGI-InSPIRE project through the European Commission's 7th
 * Framework Programme (contract # INFSO-RI-261323)
 */

package respond

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/ARGOeu/argo-web-api/utils/caches"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/logging"
	"github.com/gorilla/mux"
)

type list []interface{}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

// ConfHandler Keeps all the configuration/variables required by all the requests
type ConfHandler struct {
	Config config.Config
}

// ResponseMessage is used to construct and marshal correctly response messages
type ResponseMessage struct {
	XMLName xml.Name       `xml:"root" json:"-"`
	Status  StatusResponse `xml:"status,omitempty" json:"status,omitempty"`
	Data    interface{}    `xml:"data>result,omitempty" json:"data,omitempty"`
	Errors  interface{}    `xml:"errors>error,omitempty" json:"errors,omitempty"`
}

// StatusResponse accompanies the ResponseMessage struct to construct a response
type StatusResponse struct {
	Message string `xml:"message,omitempty" json:"message,omitempty"`
	Code    string `xml:"code,omitempty" json:"code,omitempty"`
	Details string `xml:"details,omitempty" json:"details,omitempty"`
}

// Respond will be called to answer to http requests to the PI
func (confhandler *ConfHandler) Respond(fn func(r *http.Request, cfg config.Config) (int, http.Header, []byte, error)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if r := recover(); r != nil {
				logging.HandleError(r)
			}
		}()
		code, header, output, err := fn(r, confhandler.Config)

		if code == http.StatusInternalServerError {
			log.Panic("Internal Server Error:", err)
		}

		//Add headers
		header.Set("Content-Length", fmt.Sprintf("%d", len(output)))

		for name, values := range header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		w.WriteHeader(code)
		w.Write(output)
	})

}

var acceptedContentTypes = []string{
	"application/xml",
	"application/json",
	"*/*",
}

var defaultContentType = "application/json"

// ParseAcceptHeader parses the accept header to determine the content type
func ParseAcceptHeader(r *http.Request) (string, error) {
	contentType := r.Header.Get("Accept")
	if r.Header.Get("Accept") == "" {
		return defaultContentType, nil
	}
	// contentType := httputil.NegotiateContentType(r, acceptedContentTypes, "notvalid")
	if strings.Contains(contentType, "application/json") {
		return "application/json", nil
	} else if strings.Contains(contentType, "application/xml") {
		return "application/xml", nil
	} else if strings.Contains(contentType, "*/*") {
		return "application/json", nil
	} else {
		return defaultContentType, errors.New("Not Acceptable ContentType")
	}
}

// MarshalContent marshals content using the marshaler that corresponds to the contentType parameter
func MarshalContent(doc interface{}, contentType string, prefix string, indent string) ([]byte, error) {
	var output []byte
	var err error

	if contentType == "application/xml" {
		output, err = xml.MarshalIndent(doc, prefix, indent)
	} else {
		output, err = json.MarshalIndent(doc, prefix, indent)
	}

	return output, err
}

func (confhandler *ConfHandler) walker(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	// route.Handler(route.GetHandler())
	return nil
}

// ResetCache resets the cache if it is set
func ResetCache(w http.ResponseWriter, r *http.Request, cfg config.Config) []byte {
	answer := ""
	if cfg.Server.Cache == true {
		caches.ResetCache()
		answer = "Cache Emptied"
	}
	answer = "No Caching is active"
	return []byte(answer)
}

// EndsWith creates a MatcherFunc to check that a url actual endswith a string
func EndsWith(match string) mux.MatcherFunc {
	return func(r *http.Request, route *mux.RouteMatch) bool {
		if matches, _ := regexp.MatchString(match, r.URL.Path); matches {
			return true
		}
		return false
	}
}

// UnauthorizedMessage is used to inform the user about incorrect api key and can be marshaled to xml and json
var UnauthorizedMessage = ResponseMessage{
	Status: StatusResponse{
		Message: "Unauthorized",
		Code:    "401",
		Details: "You need to provide a correct authentication token using the header 'x-api-key'",
	}}

// NotAcceptableContentType is used to inform the user about incorrect Accept header and can be marshaled to xml and json
var NotAcceptableContentType = ResponseMessage{
	Status: StatusResponse{
		Message: "Not Acceptable Content Type",
		Code:    "406",
		Details: "Accept header provided did not contain any valid content types. Acceptable content types are 'application/xml' and 'application/json'",
	}}
