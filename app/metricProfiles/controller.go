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

package metricProfiles

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
)

// ListOne handles the listing of one specific profile based on its given id
func ListOne(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Content Negotiation
	contentType, err = respond.ParseAcceptHeader(r)
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	vars := mux.Vars(r)

	if err != nil {
		code = http.StatusNotAcceptable
		output, _ = respond.MarshalContent(respond.NotAcceptableContentType, contentType, "", " ")
		return code, h, output, err
	}

	// Tenant Authentication
	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output, _ = respond.MarshalContent(respond.UnauthorizedMessage, contentType, "", " ")
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		return code, h, output, err
	}

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	filter := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, "metric_profiles", filter, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.NotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListView(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

// List the existing metric profiles for the tenant making the request
// Also there is an optional url param "name" to filter results by
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Content Negotiation
	contentType, err = respond.ParseAcceptHeader(r)
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	if err != nil {
		code = http.StatusNotAcceptable
		output, _ = respond.MarshalContent(respond.NotAcceptableContentType, contentType, "", " ")
		return code, h, output, err
	}

	urlValues := r.URL.Query()

	// Tenant Authentication
	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output, _ = respond.MarshalContent(respond.UnauthorizedMessage, contentType, "", " ")
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		return code, h, output, err
	}

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Retrieve Results from database

	var filter interface{}
	if len(urlValues["name"]) > 0 {
		filter = bson.M{"name": urlValues["name"][0]}
	} else {
		filter = nil
	}

	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, "metric_profiles", filter, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createListView(results, "Success", code) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}

//Create a new metric profile
func Create(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	contentType, err := respond.ParseAcceptHeader(r)
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	if err != nil {
		code = http.StatusNotAcceptable
		output, _ = respond.MarshalContent(respond.NotAcceptableContentType, contentType, "", " ")
		return code, h, output, err
	}

	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output, _ = respond.MarshalContent(respond.UnauthorizedMessage, contentType, "", " ")
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	incoming := MongoInterface{}

	// Try ingest request body
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	// Parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestBadJSON, contentType, "", " ")
		code = 400
		return code, h, output, err
	}

	// Generate new id
	incoming.ID = mongo.NewUUID()
	err = mongo.Insert(session, tenantDbConfig.Db, "metric_profiles", incoming)

	if err != nil {
		panic(err)
	}

	// Create view of the results
	output, err = createRefView(incoming, "Metric Profile successfully created", 201, r) //Render the results into JSON
	code = 201
	return code, h, output, err
}

//Update function to update contents of an existing metric profile
func Update(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	vars := mux.Vars(r)

	contentType, err := respond.ParseAcceptHeader(r)
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	if err != nil {
		code = http.StatusNotAcceptable
		output, _ = respond.MarshalContent(respond.NotAcceptableContentType, contentType, "", " ")
		return code, h, output, err
	}

	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output, _ = respond.MarshalContent(respond.UnauthorizedMessage, contentType, "", " ")
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		return code, h, output, err
	}

	incoming := MongoInterface{}

	// ingest body data
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	// parse body json
	if err := json.Unmarshal(body, &incoming); err != nil {
		output, _ = respond.MarshalContent(respond.BadRequestBadJSON, contentType, "", " ")
		code = 400
		return code, h, output, err
	}

	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)
	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}
	// create filter to retrieve specific profile with id
	filter := bson.M{"id": vars["ID"]}

	incoming.ID = vars["ID"]

	// Retrieve Results from database
	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, "metric_profiles", filter, "name", &results)

	if err != nil {
		panic(err)
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.NotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	// run the update query
	err = mongo.Update(session, tenantDbConfig.Db, "metric_profiles", filter, incoming)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view for response message
	output, err = createMsgView("Metric Profile successfully updated", 200) //Render the results into JSON
	code = 200
	return code, h, output, err
}

//Delete metric profile based on id
func Delete(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {

	//STANDARD DECLARATIONS START

	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	contentType := "text/xml"
	charset := "utf-8"

	//STANDARD DECLARATIONS END

	// Content Negotiation
	contentType, err = respond.ParseAcceptHeader(r)
	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))

	vars := mux.Vars(r)

	if err != nil {
		code = http.StatusNotAcceptable
		output, _ = respond.MarshalContent(respond.NotAcceptableContentType, contentType, "", " ")
		return code, h, output, err
	}

	// Tenant Authentication
	tenantDbConfig, err := authentication.AuthenticateTenant(r.Header, cfg)

	if err != nil {
		output, _ = respond.MarshalContent(respond.UnauthorizedMessage, contentType, "", " ")
		code = http.StatusUnauthorized //If wrong api key is passed we return UNAUTHORIZED http status
		return code, h, output, err
	}

	// Open session to tenant database
	session, err := mongo.OpenSession(tenantDbConfig)
	defer mongo.CloseSession(session)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	filter := bson.M{"id": vars["ID"]}

	// Retrieve Results from database
	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, "metric_profiles", filter, "name", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Check if nothing found
	if len(results) < 1 {
		output, _ = respond.MarshalContent(respond.NotFound, contentType, "", " ")
		code = 404
		return code, h, output, err
	}

	mongo.Remove(session, tenantDbConfig.Db, "metric_profiles", filter)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	// Create view of the results
	output, err = createMsgView("Metric Profile Successfully Deleted", 200) //Render the results into JSON

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	h.Set("Content-Type", fmt.Sprintf("%s; charset=%s", contentType, charset))
	return code, h, output, err
}
