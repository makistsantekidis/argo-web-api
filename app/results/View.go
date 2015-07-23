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

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func createView(results []EndpointGroupInterface, format string) ([]byte, error) {

	docRoot := &root{}

	prevReport := ""
	prevEndpointGroup := ""
	endpointGroup := &EndpointGroup{}
	report := &Report{}

	// we iterate through the results struct array
	// keeping only the value of each row

	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new report value does not match the previous report value
		//we create a new report in the xml
		if prevReport != row.Report {
			prevReport = row.Report
			report = &Report{
				Name: row.Report,
			}
			docRoot.Report = append(docRoot.Report, report)
			prevEndpointGroup = ""
		}
		//if new endpointGroup does not match the previous service value
		//we create a new endpointGroup entry in the xml
		if prevEndpointGroup != row.Name {
			prevEndpointGroup = row.Name
			endpointGroup = &EndpointGroup{
				Name:       row.Name,
				SuperGroup: row.SuperGroup,
			}
			report.EndpointGroup = append(report.EndpointGroup, endpointGroup)
		}
		//we append the new availability values
		endpointGroup.Availability = append(endpointGroup.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability),
				Unknown:      fmt.Sprintf("%g", row.Unknown)})
	}
	if strings.ToLower(format) == "json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	} else {
		return xml.MarshalIndent(docRoot, " ", "  ")
	}

}
