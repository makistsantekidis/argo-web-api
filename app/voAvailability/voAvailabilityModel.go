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

package voAvailability

import (
	"labix.org/v2/mgo/bson"
	"strconv"
	"time"
)

type Availability struct {
	Timestamp    string `xml:"timestamp,attr"`
	Availability string `xml:"availability,attr"`
	Reliability  string `xml:"reliability,attr"`
}

type Vo struct {
	Vo           string `xml:"VO,attr"`
	Availability []*Availability
}

type Profile struct {
	Name string `xml:"name,attr"`
	Vo   []*Vo
}

type Root struct {
	Profile []*Profile
}

type ApiVoAvailabilityInProfileInput struct {
	// mandatory values
	start_time           string // UTC time in W3C format
	end_time             string // UTC time in W3C format
	availability_profile string //availability profile
	granularity          string // availability period; possible values: `HOURLY`, `DAILY`, `WEEKLY`, `MONTHLY`
	// optional values
	//format    string   // default XML; possible values are: XML, JSON
	group_name []string // site name; may appear more than once
}

type ApiVoAvailabilityInProfileOutput struct {
	Date         string  "dt"
	Profile      string  "p"
	Vo           string  "v"
	Availability float64 "a"
	Reliability  float64 "r"
}

type list []interface{}

var customForm []string

func init() {
	customForm = []string{"20060102", "2006-01-02"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

func prepareFilter(input ApiVoAvailabilityInProfileInput) bson.M {

	ts, _ := time.Parse(zuluForm, input.start_time)
	te, _ := time.Parse(zuluForm, input.end_time)
	tsYMD, _ := strconv.Atoi(ts.Format(ymdForm))
	teYMD, _ := strconv.Atoi(te.Format(ymdForm))

	filter := bson.M{
		"ap": input.availability_profile,
		"dt": bson.M{"$gte": tsYMD, "$lte": teYMD},
	}

	if len(input.group_name) > 0 {
		filter["v"] = bson.M{"$in": input.group_name}
	}

	return filter
}

func Daily(input ApiVoAvailabilityInProfileInput) []bson.M {

	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 8}}}, "p": "$p", "v": "$v", "a": "$a", "r": "$r"}}},
		{"$project": bson.M{"dt": "$_id.dt", "v": "$_id.v", "p": "$_id.p", "a": "$_id.a", "r": "$_id.r"}},
		{"$sort": bson.D{{"p", 1}, {"v", 1}, {"dt", 1}}}}

	return query
}

func Monthly(input ApiVoAvailabilityInProfileInput) []bson.M {
	filter := prepareFilter(input)

	query := []bson.M{
		{"$match": filter},
		{"$group": bson.M{"_id": bson.M{"dt": bson.D{{"$substr", list{"$dt", 0, 6}}}, "p": "$p", "v": "$v"},
			"avgup": bson.M{"$avg": "$up"}, "avgu": bson.M{"$avg": "$u"}, "avgd": bson.M{"$avg": "$d"}}},
		{"$project": bson.M{"dt": "$_id.dt", "v": "$_id.v", "p": "$_id.p",
			"a": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{1.00000001, "$avgu"}}}}, 100}},
			"r": bson.M{"$multiply": list{bson.M{"$divide": list{"$avgup", bson.M{"$subtract": list{bson.M{"$subtract": list{1.00000001, "$avgu"}}, "$avgd"}}}}, 100}}}},
		{"$sort": bson.D{{"p", 1}, {"v", 1}, {"dt", 1}}}}

	return query
}
