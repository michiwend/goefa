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
	Id              int     `xml:"stopID,attr"`
	Name            string  `xml:"name,attr"`
	Locality        string  `xml:"locality,attr"`
	Lat             float64 `xml:"x,attr"`
	Lng             float64 `xml:"y,attr"`
	IsTransferStop  bool

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

type EFARoutePart struct {
    //TODO: Footpaths
    Duration            time.Duration
    Termini             []struct{
        EFARouteStop
        TimeActual      EFATime             `xml:"itdDateTime"`
        TimeTarget      EFATime             `xml:"itdDateTimeTarget"`
        Usage           string              `xml:"usage,attr"`
    }    `xml:"itdPoint"`

    MeansOfTransport    EFAMeansOfTransport `xml:"itdMeansOfTransport"`
    Stops               []*EFARouteStop     `xml:"itdStopSeq>itdPoint"`
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

func (efa *EFAProvider) Trip(origin, destination EFAStop, time time.Time, depArr string) ([]*EFARoute, error) {
    //TODO: add via routing
    params := url.Values{
        "locationServerActive":         {"1"},
        "stateless":                    {"1"},
        "itdDate":                      {time.Format("20060102")},
        "itdTime":                      {time.Format("1504")},
        "itdTripDateTimeDepArr":        {depArr},
        "nameInfo_origin":              {strconv.Itoa(origin.Id)},
        "type_origin":                  {"any"},
        "nameInfo_destination":         {strconv.Itoa(destination.Id)},
        "type_destination":             {"any"},
    }
//    if false {
//        params.Set("nameInfo_via", strconv.Itoa(via.Id))
//        params.Set("type_via", "any")
//    }

    var result tripResult

    if err := efa.postRequest(&result, params); err != nil {
        return nil, err
    }

    return result.Routes, nil
}
