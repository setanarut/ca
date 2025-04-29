package main

import (
	"fmt"
	"image"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/anthonynsimon/bild/imgio"
)

var NeumannNeighborsOffsets = [4]image.Point{
	{-1, 0}, {0, -1}, {1, 0}, {0, 1},
}

type Neighborhood int

const Neumann Neighborhood = 0
const Moore Neighborhood = 1

type RuleType int

const (
	Cyclic RuleType = iota
	Dunes
)

type Rule struct {
	Type             RuleType
	Range, Threshold int
	States           uint8
	Neighborhood     Neighborhood
}

func (r *Rule) String() string {
	var neig string
	if r.Neighborhood == Moore {
		neig = "Moore"
	} else {
		neig = "Neumann"
	}

	var ruleType string
	switch r.Type {
	case Cyclic:
		ruleType = "Cyclic"
	case Dunes:
		ruleType = "Dunes"
	}

	return fmt.Sprintf("%s_%d_%d_%d_%v", ruleType, r.Range, r.Threshold, r.States, neig)
}

// Dunes Cellular Automaton
type CCA struct {
	Grid, Buffer *image.Gray
	Neighbors    [][]int // Used for Cyclic rule
	Rnd          *rand.Rand
	Rule         Rule
	// ColorMap maps internal state values (0 to States-1) to displayable values (0-255)
	ColorMap []uint8
}

func NewAutomaton(w, h int, seed uint64, r Rule) *CCA {
	ca := &CCA{
		Grid:     image.NewGray(image.Rect(0, 0, w, h)),
		Buffer:   image.NewGray(image.Rect(0, 0, w, h)),
		Rnd:      rand.New(rand.NewPCG(seed, 0)),
		Rule:     r,
		ColorMap: make([]uint8, r.States),
	}

	// Initialize the color map to evenly spread values across 0-255 range
	for i := range ca.ColorMap {
		ca.ColorMap[i] = uint8((i * 255) / (int(r.States) - 1))
	}
	ca.initGrid()
	return ca
}

func (c *CCA) initGrid() {
	for y := range c.Grid.Rect.Dy() {
		for x := range c.Grid.Rect.Dx() {
			// Generate internal state (0 to States-1)
			state := uint8(c.Rnd.IntN(int(c.Rule.States)))
			idx := c.Grid.PixOffset(x, y)
			// Store the internal state directly, we'll map to display value during rendering
			c.Grid.Pix[idx] = state
		}
	}
}

// Get the top-left neighbor value
func (c *CCA) getTopLeftNeighbor(x, y int) uint8 {
	w, h := c.Grid.Rect.Dx(), c.Grid.Rect.Dy()
	nx := (x - 1 + w) % w
	ny := (y - 1 + h) % h
	idx := c.Grid.PixOffset(nx, ny)
	return c.Grid.Pix[idx]
}

// Get the average value of von Neumann neighbors
func (c *CCA) getVonNeumannAverage(x, y int) uint8 {
	w, h := c.Grid.Rect.Dx(), c.Grid.Rect.Dy()

	sum := 0
	count := 0
	for _, offset := range NeumannNeighborsOffsets {
		nx := (x + offset.X + w) % w
		ny := (y + offset.Y + h) % h
		idx := c.Grid.PixOffset(nx, ny)
		sum += int(c.Grid.Pix[idx])
		count++
	}
	avg := uint8(sum / count)
	return avg
}

// Helper function for absolute value
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// precomputeNeighbors prepares the neighbor lookup for cyclic rules
func (c *CCA) precomputeNeighbors() {
	w, h := c.Grid.Rect.Dx(), c.Grid.Rect.Dy()
	c.Neighbors = make([][]int, len(c.Grid.Pix))
	for y := range h {
		for x := range w {
			index := c.Grid.PixOffset(x, y)
			dirs := []int{}
			for dy := -c.Rule.Range; dy <= c.Rule.Range; dy++ {
				for dx := -c.Rule.Range; dx <= c.Rule.Range; dx++ {
					if dx == 0 && dy == 0 {
						continue
					}

					if c.Rule.Neighborhood != Moore && Abs(dx)+Abs(dy) > c.Rule.Range {
						continue
					}

					nx := (x + dx + w) % w
					ny := (y + dy + h) % h
					neighborIndex := c.Grid.PixOffset(nx, ny)
					dirs = append(dirs, neighborIndex)
				}
			}
			c.Neighbors[index] = dirs
		}
	}
}

func (c *CCA) Step() {
	switch c.Rule.Type {
	case Cyclic:
		c.cyclicStep()
	case Dunes:
		c.dunesStep()
	}
}

// cyclicStep implements the original Cyclic rule
func (c *CCA) cyclicStep() {
	w, h := c.Grid.Rect.Dx(), c.Grid.Rect.Dy()
	var wg sync.WaitGroup
	wg.Add(h)

	// Precompute neighbors if needed
	if c.Neighbors == nil {
		c.precomputeNeighbors()
	}

	// Update grid to the next generation concurrently, row by row.
	for y := range h {
		// Capture current y
		go func(y int) {
			defer wg.Done()
			for x := range w {
				idx := c.Grid.PixOffset(x, y)
				currentState := c.Grid.Pix[idx]
				targetState := (currentState + 1) % c.Rule.States
				count := 0

				for _, nIdx := range c.Neighbors[idx] {
					if c.Grid.Pix[nIdx] == targetState {
						count++
						if count >= c.Rule.Threshold {
							break
						}
					}
				}

				if count >= c.Rule.Threshold {
					c.Buffer.Pix[idx] = targetState
				} else {
					c.Buffer.Pix[idx] = currentState
				}
			}
		}(y)
	}

	wg.Wait()

	// Swap grid and buffer
	c.Grid, c.Buffer = c.Buffer, c.Grid
} // dunesStep implements the Dunes rule
func (c *CCA) dunesStep() {
	w, h := c.Grid.Rect.Dx(), c.Grid.Rect.Dy()
	var wg sync.WaitGroup
	wg.Add(h)

	// Update grid to the next generation concurrently, row by row
	for y := range h {
		// Capture current y
		go func(y int) {
			defer wg.Done()
			for x := range w {
				idx := c.Grid.PixOffset(x, y)
				currentState := c.Grid.Pix[idx]

				// Get top-left neighbor value
				topLeftValue := c.getTopLeftNeighbor(x, y)

				var newState uint8
				// If top-left neighbor value > N/2
				if topLeftValue > c.Rule.States/2 {
					// Increase current value by 1 (mod N)
					newState = (currentState + 1) % c.Rule.States
				} else {
					// Get average of von Neumann neighbors
					avgValue := c.getVonNeumannAverage(x, y)
					// Set to successor of average (mod N)
					newState = (avgValue + 1) % c.Rule.States
				}

				c.Buffer.Pix[idx] = newState
			}
		}(y)
	}

	wg.Wait()

	// Swap grid and buffer
	c.Grid, c.Buffer = c.Buffer, c.Grid
}

// RenderedImage returns a new image.Gray with color-mapped pixel values
func (c *CCA) RenderedImage() *image.Gray {
	w, h := c.Buffer.Rect.Dx(), c.Buffer.Rect.Dy()
	img := image.NewGray(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			idx := c.Buffer.PixOffset(x, y)
			state := c.Buffer.Pix[idx]
			img.Pix[idx] = c.ColorMap[state]
		}
	}
	return img
}

func main() {
	// otomat := NewAutomaton(512, 256, 33, DunesRule(32))
	otomat := NewAutomaton(512, 256, 33, CyclicRule(1, 1, 24, Moore))

	t := time.Now()
	for range 500 {
		otomat.Step()
	}
	imgio.Save("output/çokluKural.png", otomat.RenderedImage(), imgio.PNGEncoder())
	fmt.Println("Simulation time taken:", time.Since(t).Seconds(), "seconds")
}

func CyclicRule(r, t int, s uint8, n Neighborhood) Rule {
	return Rule{
		Type:         Cyclic,
		Range:        r,
		Threshold:    t,
		States:       s,
		Neighborhood: n,
	}
}
func DunesRule(states uint8) Rule {
	return Rule{
		Type:   Dunes,
		States: states,
	}
}
