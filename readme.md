# EquiEasy
Read Equibase Horse Racing PDF sheets in to data structures.

## Install:
`go get -u github.com/hektorinho/equieasy`

## Example:
Gather all data from an Equibase PDF:
```golang
package main

import (
    "log"

    "github.com/hektorinho/equieasy"
)

const (
	file = "data/eqbPDFChartPlus.pdf"
)

func main() {
	r, err := pdf.Open(testFile)
	if err != nil {
		log.Panicln(err)
	}
	p, err := GetValidPages(testFile, r)
	if err != nil {
		log.Panicln(err)
	}

	for _, page := range p.Pages {
		race, err := NewRacePage(page)
		if err != nil {
		    log.Panicln(err)
		}

		if err := DoSomethingWithData(race); err != nil {
			log.Panicln(err)
		}
	}
}

func DoSomethingWithData(race RacePage) error {
        //TODO:
        return nil
}
