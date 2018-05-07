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
	"encoding/xml"
	"time"
    "strings"
)

// EFATime implements UnmarshalXML to support unmarshalling EFAs XML DateTime
// type directly into a time.Time compatible type
type EFATime struct {
	*time.Time
}

func (t *EFATime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	type efaDateTime struct {
		Date struct {
			Day   int `xml:"day,attr"`
			Month int `xml:"month,attr"`
			Year  int `xml:"year,attr"`
		} `xml:"itdDate"`

		Time struct {
			Hour   int `xml:"hour,attr"`
			Minute int `xml:"minute,attr"`
		} `xml:"itdTime"`
	}

	var content efaDateTime

	if err := d.DecodeElement(&content, &start); err != nil {
		return err
	}

	loc, err := time.LoadLocation("Local")

	if err != nil {
		return err
	}

	tmp := time.Date(content.Date.Year,
		time.Month(content.Date.Month),
		content.Date.Day,
		content.Time.Hour,
		content.Time.Minute,
		0,
		0,
		loc)

	t.Time = &tmp

	return nil
}

// Durations are stored in *three* different ways in the XML:
// 1. In one Attribute, just in minutes
// 2. In one Attribute, as "hh:mm"
// 3. In multiple Attributes, as "timeMinute", "timeHour" (needs more investigating)
// 
// If EFADuration is unmarshalled on a single xml attribute, 1 or 2 is used, 
// if unmarshalled on an xml element, case 3 is used.
type EFADuration struct {
    *time.Duration
}



func (dur *EFADuration) UnmarshalXMLAttr(attr xml.Attr) error {
    var duration time.Duration
    var err error

    if strings.Contains(attr.Value, ":"){
        timestring := strings.Replace(attr.Value, ":", "h", 1) + "m"
        duration, err = time.ParseDuration(timestring)
    } else {
        timestring := attr.Value + "m"
        duration, err = time.ParseDuration(timestring)
    }

    if err != nil {
        return err
    }
    dur.Duration = &duration
    return nil
}
