package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/michiwend/goefa"
)

func main() {
    provider, _ := goefa.ProviderFromJson("mvv")

    _, stops, err := provider.FindStop("Innsbrucker Ring")

    if err != nil {
        log.Panic(err)
    }

    from := stops[0]

    _, stops, err = provider.FindStop("Hauptbahnhof")

    if err != nil {
        log.Panic(err)
    }

    to := stops[0]

    routes, err := provider.Trip(*from, *to, time.Now(), "dep")
    
    if err != nil {
        log.Panic(err)
    }

    fmt.Println(routes)

    for _, route := range routes {
        fmt.Printf("Parts: %+v, %+v\n", route.RouteParts, route.PublicDuration)
        for _, routePart := range route.RouteParts {
            for _, ter := range routePart.Termini {
                fmt.Printf("Termini: %+v\n", ter.Name)
            }
            fmt.Printf("MOT: %+v\n", routePart.MeansOfTransport)

            for _, stop := range routePart.Stops {
                fmt.Printf("\t%+v\n", stop)
            }
        }
    }
}
