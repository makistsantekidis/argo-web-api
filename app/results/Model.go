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

import "encoding/xml"

type list []interface{}

var customForm []string

func init() {
	customForm = []string{"20060102", "2006-01-02T15:04:05Z"} //{"Format that is returned by the database" , "Format that will be used in the generated report"}
}

const zuluForm = "2006-01-02T15:04:05Z"
const ymdForm = "20060102"

type endpointGroupResultQuery struct {
	Name        string `bson:"name"`
	Granularity string `bson:"-"`
	Format      string `bson:"-"`
	StartTime   string `bson:"start_time"` // UTC time in W3C format
	EndTime     string `bson:"end_time"`   // UTC time in W3C format
	Report      string `bson:"report"`
	Group       string `bson:"supergroup"`
}

type EndpointGroupInterface struct {
	Name         string  `bson:"name"`
	Report       string  `bson:"report"`
	Date         string  `bson:"date"`
	Type         string  `bson:"type"`
	Up           float64 `bson:"uptime"`
	Down         float64 `bson:"downtime"`
	Unknown      float64 `bson:"unknown"`
	Availability float64 `bson:"availability"`
	Reliability  float64 `bson:"reliability"`
	Weights      string  `bson:"weights"`
	SuperGroup   string  `bson:"supergroup"`
}

//Availability struct for formating xml/json
type Availability struct {
	XMLName      xml.Name `xml:"Availability" json:"-"`
	Timestamp    string   `xml:"timestamp,attr" json:"timestamp"`
	Availability string   `xml:"availability,attr" json:"availability"`
	Reliability  string   `xml:"reliability,attr" json:"reliability"`
	Unknown      string   `xml:"unknown,attr" json:"unknown"`
}

//EndpointGroup struct for formating xml/json
type EndpointGroup struct {
	XMLName      xml.Name `xml:"EndpointGroup" json:"-"`
	Name         string   `xml:"name,attr" json:"name"`
	SuperGroup   string   `xml:"group,attr" json:"group"`
	Availability []interface{}
}

//Report struct for formating xml/json
type Report struct {
	XMLName       xml.Name `xml:"Report" json:"-"`
	Name          string   `xml:"name,attr" json:"name"`
	EndpointGroup []interface{}
}

type root struct {
	XMLName xml.Name `xml:"root" json:"-"`
	Report  []interface{}
}
