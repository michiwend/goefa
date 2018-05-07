/*
 * Copyright (C) 2014 Michael Wendland
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Authors:
 *   Michael Wendland <michael@michiwend.com>
 */

package goefa

import (
	"net/url"
	"strconv"
	"time"
)

// EFADepartureArrival represents either an arrival or a departure and is not
// directly used but embedded by EFAArrival and EFADeparture.
type EFADepartureArrival struct {
	Area         int    `xml:"area,attr"`
	Countdown    int    `xml:"countdown,attr"`
	MapName      string `xml:"mapName,attr"`
	Platform     string `xml:"platform,attr"`
	PlatformName string `xml:"platformName,attr"`
	StopID       int    `xml:"stopID,attr"`
	StopName     string `xml:"stopName,attr"`
	Lat          int64  `xml:"x,attr"`
	Lng          int64  `xml:"y,attr"`

	DateTime    EFATime        `xml:"itdDateTime"`
	ServingLine EFAServingLine `xml:"itdServingLine"`
}

// EFAArrival represents a single arrival for a specific stop.
type EFAArrival struct {
	EFADepartureArrival `xml:"itdArrival"`
}

// EFADeparture represents a single departure for a specific stop.
type EFADeparture struct {
	EFADepartureArrival `xml:"itdDeparture"`
}

type odv struct {
    OdvPlace struct {
    }
    OdvName struct {
        State string `xml:"state,attr"`
    } `xml:"itdOdvName"`
}

type departureMonitorResult struct {
	EFAResponse
	Odv odv `xml:"itdDepartureMonitorRequest>itdOdv"`
	Lines		[]*EFAServingLine	`xml:"itdDepartureMonitorRequest>itdServingLines>itdServingLine"`
	Departures	[]*EFADeparture		`xml:"itdDepartureMonitorRequest>itdDepartureList>itdDeparture"`
}

func (d *departureMonitorResult) endpoint() string {
	return "XML_DM_REQUEST"
}

type DepartureRequest struct {
	Params				*url.Values
	StopId				int
	Time				time.Time
	Results				int
	Lines				[]*EFAServingLine
}

func (dr *DepartureRequest) getDefaultParams() url.Values {
	params := url.Values{
		"type_dm":				{"any"},
		"locationServerActive": {"1"},
		"mode":					{"direct"},
		"stateless":			{"1"},
	}
	return params
}

func (dr *DepartureRequest) GetParams() url.Values {
	var params url.Values
	if dr.Params == nil {
		params = dr.getDefaultParams()
	} else {
		params = *dr.Params
	}

	params.Set("name_dm", strconv.Itoa(dr.StopId))
	params.Set("itdDate", dr.Time.Format("20060102"))
	params.Set("itdTime", dr.Time.Format("1504"))
	params.Set("limit", strconv.Itoa(dr.Results))

	for _, line := range dr.Lines {
		params.Add("line", line.Stateless)
	}
	return params
}

// Departures performs a stateless dm_request for the corresponding stopID and
// returns an array of EFADepartures. Use time.Now() as the second argument in
// order to get the very next departures. The third argument determines how
// many results will be returned by EFA.
func (efa *EFAProvider) Departures(stopID int, due time.Time, results int) ([]*EFADeparture, error) {
    return efa.DeparturesForLines(stopID, due, results, nil)
}

func (efa *EFAProvider) DeparturesForLines(stopID int, due time.Time, results int, lines []*EFAServingLine) ([]*EFADeparture, error) {
	var rt string

	if efa.EnableRealtime {
		rt = "1"
	} else {
		rt = "0"
	}

	req := DepartureRequest{
		StopId:		stopID,
		Time:		due,
		Results:	results,
		Lines:      lines,
	}

	params := req.GetParams()
	params.Set("useRealtime", rt)

	var result departureMonitorResult

	if err := efa.postRequest(&result, params); err != nil {
		return nil, err
	}

	return result.Departures, nil
}

// Lines performs a stateless dm_request for the corresponding stopID and
// returns an array of EFAServingLines.
func (efa *EFAProvider) Lines(stopID int) ([]*EFAServingLine, error) {
	req := DepartureRequest{
		StopId:		stopID,
		Time:		time.Now(),
	}

	params := req.GetParams()

	var result departureMonitorResult

	if err := efa.postRequest(&result, params); err != nil {
		return nil, err
	}

	return result.Lines, nil
}
