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
)


type tripResult struct {
   EFAResponse
   Odv []struct {
    odv
    Usage string `xml:"usage,attr"`
   } `xml:"itdTripRequest>itdOdv"`

   Xml string `xml:",innerxml"`
}

func (t *tripResult) endpoint() string {
    return "XML_TRIP_REQUEST2"
}

func (efa *EFAProvider) Trip(origin, destination EFAStop, time time.Time, depArr string) (*tripResult, error) {
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

    return &result, nil
}
