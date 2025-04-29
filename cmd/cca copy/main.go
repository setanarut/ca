package main

import (
	"ca/cca"
	"fmt"
	"image"
	"math"
	"strconv"
	"time"

	"github.com/setanarut/apng"
)

var Rule_1_1_16_Neumann = cca.Rule{
	Range:        1,
	Threshold:    1,
	States:       16,
	Neighborhood: cca.Neumann,
}
var Rule_1_1_24_Moore = cca.Rule{
	Range:        1,
	Threshold:    1,
	States:       24,
	Neighborhood: cca.Moore,
}

func main() {
	rule := cca.Rule{
		Range:        2,
		Threshold:    2,
		States:       23,
		Neighborhood: cca.Moore,
	}

	maxGeneration := 600
	sim := cca.NewCyclicAutomaton(800, int(math.Round(800.0*0.618)), 2, rule)
	// seri := make([]uint8, 0)

	t := time.Now()
	for range maxGeneration {
		// seri = append(seri, sim.Grid.Pix[0])
		sim.Step()

	}
	fmt.Println(maxGeneration, "Frames time taken:", time.Since(t).Seconds(), "seconds")

	// n, s, ln := cca.FindLastLongestConsecutivePattern(seri)

	frames := make([]image.Image, 0)

	start := sim.Grid.Pix[0]
	loopEnd := 0
	for i := range 100 {
		frames = append(frames, cca.GetFrame(sim))
		sim.Step()
		if start == sim.Grid.Pix[0] {
			loopEnd = i + maxGeneration
			fmt.Println(i, "Frames loop")
			break
		}
	}
	// fmt.Println(n, s, ln)
	prefix := sim.Rule.String() + "_" + strconv.Itoa(maxGeneration) + "-" + strconv.Itoa(loopEnd)
	fileName := "output/out_" + prefix + ".png"
	apng.Save(fileName, frames, 3)
	// imgio.Save("out/out_"+strconv.Itoa(total)+".png", cca.GetFrame(sim), imgio.PNGEncoder())
	// imgio.Save("frame"+strconv.Itoa(total)+".png", sim.Image(), imgio.PNGEncoder())
}
