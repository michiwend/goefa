/*
 *    Copyright (C) 2014 Michael Wendland
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Michael Wendland <michiwend@michiwend.com>
 */

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/leftshift/goefa"
)

func main() {

	pname := flag.String("provider", "avv", "Short name for the EFA Provider")
	query := flag.String("stop", "Koenigsplatz", "The stop name to search for")
	flag.Parse()

	myprovider, err := goefa.ProviderFromJson(*pname)

	if err != nil {
		log.Println(err)
		return
	}

	idtfd, stops, err := myprovider.FindStop(*query)

	if err != nil {
		log.Println(err)
		return
	}

	var mystop *goefa.EFAStop

	if err != nil {
		log.Println(err)
		return
	}

	if idtfd == false {
		fmt.Println("Two or more stops where matched:")
		for i, stop := range stops {

			fmt.Printf("%2d - %s (%s)\n", i, stop.Name, stop.Locality)
		}

		fmt.Print("Choose one: ")
		var i int
		_, err := fmt.Scanf("%d", &i)

		if err != nil {
			log.Println(err)
			return
		}

		if i > len(stops) {
			log.Println(errors.New("Index out of range."))
			return
		}

		mystop = stops[i]

	} else {
		mystop = stops[0]
		fmt.Println("Stop identified: " + mystop.Name)
	}
    fmt.Println("Lines serving this station:")

	lines, err := mystop.Lines()

	if err != nil {
		log.Println(err)
		return
	}

	for _, line := range lines {

		fmt.Printf("%17s %-5s --> %s\n",
			line.MotType.String(),
			line.Number,
			line.Direction)

	}

}
