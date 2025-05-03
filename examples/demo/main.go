package main

import (
	"ca"
	"log"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/mazznoer/colorgrad"
)

func main() {

	otomat := ca.NewAutomaton(854, 480, 33, ca.DunesRule(32))
	otomat.FillWithRandomStates()

	for range 30 * 20 {
		otomat.Step()
	}

	err := imgio.Save(
		"examples/demo/o.png",
		otomat.GradientMap(colorgrad.Rainbow()),
		imgio.PNGEncoder(),
	)
	if err != nil {
		log.Fatal(err)
	}
}
