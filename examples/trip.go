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
