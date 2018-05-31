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
 *   Adrian Schneider <github@ardy.io>
 */

package goefa

import (
    "time"
    "strconv"
    "net/url"
    "encoding/xml"
)


type EFARouteStop struct {
    //FIXME: This is basically the same as EFAStop with slightly different attr names
    //(Though it actually can represent other points, e.g. within stations)
	Id              int     `xml:"stopID,attr"`
	Name            string  `xml:"name,attr"`
	Locality        string  `xml:"locality,attr"`
	Lat             float64 `xml:"x,attr"`
	Lng             float64 `xml:"y,attr"`
    Area            int     `xml:"area,attr"`

    Platform struct {
        Id          string  `xml:"platform"`
        Name        string  `xml:"platformName"`
    }

    Times []EFATime         `xml:"itdDateTime"`

	Provider *EFAProvider
}

type EFAMeansOfTransport struct {
    Name            string              `xml:"name,attr"`
    Shortname       string              `xml:"shortname,attr"`
    Symbol          string              `xml:"symbol,attr"`
    Type            EFAMotType          `xml:"motType,attr"`
    ProductName     string              `xml:"productName,attr"`
    Destination     string              `xml:"destination,attr"`
    DestId          int                 `xml:"destID,attr"`
    Network         string              `xml:"network,attr"`

	ROP             int                 `xml:"ROP,attr"`
	STT             int                 `xml:"STT,attr"`
	TTB             int                 `xml:"TTB,attr"`

    Description     string              `xml:"itdRouteDescText"`
}

type EFAFootpathElem struct {
    Description     string              `xml:"description,attr"`
    Type            string              `xml:"type,attr"`
    VertDirection   string              `xml:"level,attr"`
    Points          []*EFARouteStop     `xml:"itdPoint"`
}

type EFAFootpathInfo struct {
    Position        string              `xml:"position,attr"`
    Duration        EFADuration         `xml:"duration,attr"`
    Elements        []*EFAFootpathElem  `xml:"itdFootPathElem"`
}

type EFARoutePart struct {
    //TODO: Footpath description/coordinates
    Duration            time.Duration
    Termini             []struct{
        EFARouteStop
        TimeActual      EFATime             `xml:"itdDateTime"`
        TimeTarget      EFATime             `xml:"itdDateTimeTarget"`
        Usage           string              `xml:"usage,attr"`
    }    `xml:"itdPoint"`

    MeansOfTransport    EFAMeansOfTransport `xml:"itdMeansOfTransport"`
    Stops               []*EFARouteStop     `xml:"itdStopSeq>itdPoint"`
    Footpath            struct {
        // Contains info for within station, e.g. stairs, ramps, escalators
        EFAFootpathInfo                     `xml:"itdFootPathInfo"`
    }
}

func (rp *EFARoutePart) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

    // This is really shitty code but I couldn't find a better way of doing it :/
    type tmp struct {
        Minutes         int                 `xml:"timeMinute,attr"`
        Hours           int                 `xml:"timeHour"`
        Termini             []struct{
            EFARouteStop
            TimeActual      EFATime             `xml:"itdDateTime"`
            TimeTarget      EFATime             `xml:"itdDateTimeTarget"`
            Usage           string              `xml:"usage,attr"`
        }    `xml:"itdPoint"`

        MeansOfTransport    EFAMeansOfTransport `xml:"itdMeansOfTransport"`
        Stops               []*EFARouteStop     `xml:"itdStopSeq>itdPoint"`
        Footpath            struct {
            // Contains info for within station, e.g. stairs, ramps, escalators
            EFAFootpathInfo                     `xml:"itdFootPathInfo"`
        }
    }

    var content tmp

    if err := d.DecodeElement(&content, &start); err != nil {
        return err
    }

    ns := content.Minutes * 60000000000 +
          content.Hours * 3600000000000

    dur := time.Duration(ns)

    rp.Duration = dur
    rp.Termini = content.Termini
    rp.MeansOfTransport = content.MeansOfTransport
    rp.Stops = content.Stops
    rp.Footpath = content.Footpath

    return nil
}

type EFARoute struct {
    DurationPublic      EFADuration         `xml:"publicDuration,attr"`
    DurationIndividual  EFADuration         `xml:"individualDuration,attr"`
    DurationVehicle     EFADuration         `xml:"vehicleTime,attr"`

    RouteParts          []*EFARoutePart  `xml:"itdPartialRouteList>itdPartialRoute"`
}

type tripResult struct {
   EFAResponse
   Odv      []struct {
        odv
        Usage   string  `xml:"usage,attr"`
   } `xml:"itdTripRequest>itdOdv"`

   Routes   []*EFARoute  `xml:"itdTripRequest>itdItinerary>itdRouteList>itdRoute"`

   Xml string `xml:",innerxml"`
}

func (t *tripResult) endpoint() string {
    return "XML_TRIP_REQUEST2"
}

type TripRequest struct {
    Params          *url.Values
    Time            time.Time
    DepArr          string
    Origin          EFAStop
    Destination     EFAStop
    Via             EFAStop
    IncludeMOT      []EFAMotType
    Results         int
}

func (tr *TripRequest) getDefaultParams() url.Values {
    params := url.Values{
        "locationServerActive":         {"1"},
        "stateless":                    {"1"},
        "type_origin":                  {"any"},
        "type_destination":             {"any"},
        "type_via":                     {"any"},
    }
    return params
}

func (tr *TripRequest) GetParams() url.Values {
    var params url.Values
    if tr.Params == nil {
        params = tr.getDefaultParams()
    } else {
        params = *tr.Params
    }

    params.Set("itdDate", tr.Time.Format("20060102"))
    params.Set("itdTime", tr.Time.Format("1504"))
    params.Set("itdTripDateTimeDepArr", tr.DepArr)
    params.Set("nameInfo_origin", strconv.Itoa(tr.Origin.Id))
    params.Set("nameInfo_destination", strconv.Itoa(tr.Destination.Id))
    params.Set("calcNumberOfTrips", strconv.Itoa(tr.Results))
    if tr.Via.Id != 0 {
        params.Set("name_via", strconv.Itoa(tr.Via.Id))
    }
    if len(tr.IncludeMOT) > 0 {
        //params.Set("includedMeans", "1")
        for _, mot := range tr.IncludeMOT {
            params.Set("inclMOT_" + strconv.Itoa(int(mot)), "on")
        }
    }
    return params
}

func (efa *EFAProvider) DoTripRequest(req *TripRequest) (*tripResult, error) {
    var result tripResult
    params := req.GetParams()
    if err := efa.postRequest(&result, params); err != nil {
        return nil, err
    }
    return &result, nil
}

func (efa *EFAProvider) Trip(origin, destination EFAStop, time time.Time, depArr string, results int) ([]*EFARoute, error) {
    //TODO: add mobility and routing preferences
    req := TripRequest{
         Origin:        origin,
         Destination:   destination,
         Time:          time,
         DepArr:        depArr,
         Results:       results,
    }

    res, err := efa.DoTripRequest(&req)
    if err != nil {
        return nil, err
    }
    return res.Routes, nil
}

func (efa *EFAProvider) TripVia(origin, via, destination EFAStop, time time.Time, depArr string, results int) ([]*EFARoute, error) {
    req := TripRequest{
         Origin:        origin,
         Via:           via,
         Destination:   destination,
         Time:          time,
         DepArr:        depArr,
         Results:       results,
    }

    res, err := efa.DoTripRequest(&req)
    if err != nil {
        return nil, err
    }
    return res.Routes, nil
}

func (efa *EFAProvider) TripUsingMot(origin, destination EFAStop, time time.Time, depArr string, mots []EFAMotType, results int) ([]*EFARoute, error) {
    //TODO: add mobility and routing preferences
    req := TripRequest{
         Origin:        origin,
         Destination:   destination,
         Time:          time,
         DepArr:        depArr,
         IncludeMOT:    mots,
         Results:       results,
    }

    res, err := efa.DoTripRequest(&req)
    if err != nil {
        return nil, err
    }
    return res.Routes, nil
}
