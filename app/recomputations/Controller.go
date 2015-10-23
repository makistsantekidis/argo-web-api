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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ARGOeu/argo-web-api/respond"
	"github.com/ARGOeu/argo-web-api/utils/authentication"
	"github.com/ARGOeu/argo-web-api/utils/config"
	"github.com/ARGOeu/argo-web-api/utils/mongo"
	"github.com/gorilla/mux"
)

var recomputationsColl = "recomputations"

// List existing recomputations
func List(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	// contentType := "application/json"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	urlValues := r.URL.Query()
	contentType := r.Header.Get("Accept")

	contentType, err = respond.ParseAcceptHeader(r)
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

	filter := IncomingRecomputation{
		StartTime: urlValues.Get("start_time"),
		EndTime:   urlValues.Get("end_time"),
		Reason:    urlValues.Get("reason"),
		Report:    urlValues.Get("report"),
	}

	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	results := []MongoInterface{}
	err = mongo.Find(session, tenantDbConfig.Db, recomputationsColl, filter, "timestamp", &results)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createListView(results, contentType)

	return code, h, output, err

}

// ListOne lists a single recomputation according to the given id
func ListOne(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusOK
	h := http.Header{}
	output := []byte("")
	err := error(nil)
	// contentType := "application/json"
	charset := "utf-8"
	//STANDARD DECLARATIONS END

	// urlValues := r.URL.Query()
	contentType := r.Header.Get("Accept")
	vars := mux.Vars(r)

	contentType, err = respond.ParseAcceptHeader(r)
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

	filter := IncomingRecomputation{
		ID: vars["ID"],
	}
	session, err := mongo.OpenSession(tenantDbConfig)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	result := MongoInterface{}
	err = mongo.FindOne(session, tenantDbConfig.Db, recomputationsColl, filter, &result)

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	output, err = createListView(result, contentType)

	return code, h, output, err

}

// SubmitRecomputation insert a new pending recomputation in the tenants database
func SubmitRecomputation(r *http.Request, cfg config.Config) (int, http.Header, []byte, error) {
	//STANDARD DECLARATIONS START
	code := http.StatusAccepted
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

	if err != nil {
		code = http.StatusInternalServerError
		return code, h, output, err
	}

	var recompSubmission IncomingRecomputation
	// urlValues := r.URL.Query()

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, cfg.Server.ReqSizeLimit))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &recompSubmission); err != nil {
		code = 422 // unprocessable entity
		output = []byte("Unprocessable JSON")
		return code, h, output, err
	}
	now := time.Now()
	recomputation := MongoInterface{
		ID:             mongo.NewUUID(),
		RequesterName:  tenantDbConfig.User,
		RequesterEmail: tenantDbConfig.Email,
		StartTime:      recompSubmission.StartTime,
		EndTime:        recompSubmission.EndTime,
		Reason:         recompSubmission.Reason,
		Report:         recompSubmission.Report,
		Exclude:        recompSubmission.Exclude,
		Timestamp:      now.Format("2006-01-02 15:04:05"),
		Status:         "pending",
	}

	err = mongo.Insert(session, tenantDbConfig.Db, recomputationsColl, recomputation)

	if err != nil {
		panic(err)
	}

	output, err = createSubmitView(recomputation, contentType, r)
	return code, h, output, err
}
