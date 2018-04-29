package main

import (
    "fmt"
    "log"
    "flag"
    "time"
    
    "github.com/michiwend/goefa"
)

func main() {

	pname := flag.String("provider", "mvv", "Short name for the EFA Provider")
    start := flag.String("start", "Moosach", "Start station for Trip")
    destination := flag.String("destination", "Obersendling", "Destination for Trip")
    flag.Parse()

    provider, err := goefa.ProviderFromJson(*pname)

    if err != nil {
        log.Panic(err)
    }

    _, stops, err := provider.FindStop(*start)

    if err != nil {
        log.Panic(err)
    }

    from := stops[0]

    _, stops, err = provider.FindStop(*destination)

    if err != nil {
        log.Panic(err)
    }

    to := stops[0]

    routes, err := provider.Trip(*from, *to, time.Now(), "dep")
    
    if err != nil {
        log.Panic(err)
    }

    fmt.Printf("Found %v routes:\n", len(routes))

    for i, route := range routes {
        fmt.Printf("Route %v, takes %v\n", i, route.DurationPublic)

        for _, routePart := range route.RouteParts {
            f, t := routePart.Termini[0], routePart.Termini[1]
            mot := routePart.MeansOfTransport
            fmt.Printf("\tFrom %v, %v to %v, %v\n", f.Name, f.TimeActual, t.Name, t.TimeActual)
            fmt.Printf("\tUsing %v to %v\n", mot.Shortname, mot.Destination)

            for _, stop := range routePart.Stops {
                t := len(stop.Times) - 1
                fmt.Printf("\t\t%v\t%v\n", stop.Times[t], stop.Name)
            }
        }
    }
}
