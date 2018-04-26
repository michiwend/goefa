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

    fmt.Println(routes.Xml)
}
