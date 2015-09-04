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
	"encoding/xml"

	"gopkg.in/mgo.v2/bson"
)

type root struct {
	MetricProfiles []MongoInterface
}

// MongoInterface to retrieve and insert metricProfiles in mongo
type MongoInterface struct {
	ID       bson.ObjectId `bson:"_id,omitempty" xml:"-"`
	OutID    string        `bson:"-" xml:"id,attr"`
	Name     string        `bson:"name" xml:"name,attr" json:"name"`
	Services []Service     `bson:"services" xml:"services" json:"services"`
}

// Service struct to represent services with their metrics
type Service struct {
	Service string   `bson:"service" xml:"service,attr" json:"service"`
	Metrics []string `bson:"metrics" xml:"metrics" json:"metrics"`
}

//Profiles to preserve compatibility with previous poem request
type Profiles struct {
	XMLName xml.Name `xml:"root"`
	Poems   []Poem   `xml:"Poem"`
}

// Poem to preserve compatibility with previous poem request
type Poem struct {
	Profile string `xml:"profile,attr" bson:"name"`
}
