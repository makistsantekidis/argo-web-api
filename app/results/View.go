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

package results

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func createServiceFlavorResultView(results []ServiceFlavorInterface, report ReportInterface, format string) ([]byte, error) {

	docRoot := &root{}

	prevSuperGroup := ""
	prevServiceFlavor := ""
	serviceFlavor := &ServiceFlavor{}
	superGroup := &SuperGroup{}

	// we iterate through the results struct array
	// keeping only the value of each row

	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new superGroup value does not match the previous superGroup value
		//we create a new superGroup in the xml
		if prevSuperGroup != row.SuperGroup {
			prevSuperGroup = row.SuperGroup
			superGroup = &SuperGroup{
				Name: row.SuperGroup,
				Type: report.EndpointGroupType, // Endpoint groups are parents of SFs
			}
			docRoot.Result = append(docRoot.Result, superGroup)
			prevServiceFlavor = ""
		}
		//if new endpointGroup does not match the previous service value
		//we create a new endpointGroup entry in the xml
		if prevServiceFlavor != row.Name {
			prevServiceFlavor = row.Name
			serviceFlavor = &ServiceFlavor{
				Name: row.Name,
				Type: fmt.Sprintf("service"),
			}
			superGroup.ServiceFlavor = append(superGroup.ServiceFlavor, serviceFlavor)
		}
		//we append the new availability values
		serviceFlavor.Availability = append(serviceFlavor.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability),
				Unknown:      fmt.Sprintf("%g", row.Unknown),
				Uptime:       fmt.Sprintf("%g", row.Up),
				Downtime:     fmt.Sprintf("%g", row.Down),
			})
	}
	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	} else {
		return xml.MarshalIndent(docRoot, " ", "  ")
	}

}

func createEndpointGroupResultView(results []EndpointGroupInterface, report ReportInterface, format string) ([]byte, error) {

	docRoot := &root{}

	prevSuperGroup := ""
	prevEndpointGroup := ""
	endpointGroup := &EndpointGroup{}
	superGroup := &SuperGroup{}

	// we iterate through the results struct array
	// keeping only the value of each row

	for _, row := range results {
		timestamp, _ := time.Parse(customForm[0], fmt.Sprint(row.Date))
		//if new superGroup value does not match the previous superGroup value
		//we create a new superGroup in the xml
		if prevSuperGroup != row.SuperGroup {
			prevSuperGroup = row.SuperGroup
			superGroup = &SuperGroup{
				Name: row.SuperGroup,
				Type: report.SuperGroupType,
			}
			docRoot.Result = append(docRoot.Result, superGroup)
			prevEndpointGroup = ""
		}
		//if new endpointGroup does not match the previous service value
		//we create a new endpointGroup entry in the xml
		if prevEndpointGroup != row.Name {
			prevEndpointGroup = row.Name
			endpointGroup = &EndpointGroup{
				Name: row.Name,
				// SuperGroup: row.SuperGroup,
				Type: report.EndpointGroupType,
			}
			superGroup.EndpointGroup = append(superGroup.EndpointGroup, endpointGroup)
		}
		//we append the new availability values
		endpointGroup.Availability = append(endpointGroup.Availability,
			&Availability{
				Timestamp:    timestamp.Format(customForm[1]),
				Availability: fmt.Sprintf("%g", row.Availability),
				Reliability:  fmt.Sprintf("%g", row.Reliability),
				Unknown:      fmt.Sprintf("%g", row.Unknown),
				Uptime:       fmt.Sprintf("%g", row.Up),
				Downtime:     fmt.Sprintf("%g", row.Down),
			})
	}
	if strings.ToLower(format) == "application/json" {
		return json.MarshalIndent(docRoot, " ", "  ")
	} else {
		return xml.MarshalIndent(docRoot, " ", "  ")
	}

}
