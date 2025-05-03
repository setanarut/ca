package ca

import (
	"fmt"
	"image"
	"image/color"
	"math/rand/v2"
	"slices"
	"sync"

	"github.com/mazznoer/colorgrad"
)

var MooreNeighborsOffsets = [8]image.Point{
	{-1, -1}, {0, -1}, {1, -1}, {-1, 0},
	{1, 0}, {-1, 1}, {0, 1}, {1, 1}}

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
	LifeLike
)

type Rule struct {
	Type             RuleType
	Range, Threshold int
	States           uint8
	Neighborhood     Neighborhood
	Burn             []int // Burn
	Survive          []int // Survive
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
	Current, Next      *image.Gray
	Neighbors          [][]int // Used for Cyclic rule
	Rnd                *rand.Rand
	Rule               Rule
	LifeLikeAgeEnabled bool
	LifeLikeMaxAge     int
}

func NewAutomaton(w, h int, seed uint64, r Rule) *CCA {
	ca := &CCA{
		Current:        image.NewGray(image.Rect(0, 0, w, h)),
		Next:           image.NewGray(image.Rect(0, 0, w, h)),
		Rnd:            rand.New(rand.NewPCG(seed, 0)),
		Rule:           r,
		LifeLikeMaxAge: 64,
	}

	// ca.Init()
	return ca
}

func (c *CCA) Init() {
	switch c.Rule.Type {
	case Cyclic:
		c.FillWithRandomStates()
	case Dunes:
		c.FillWithRandomStates()
	case LifeLike:
		c.FillWithRandom_0_255()
	}
}

func (c *CCA) FillWithRandomStates() {
	for y := range c.Current.Rect.Dy() {
		for x := range c.Current.Rect.Dx() {
			c.Current.SetGray(x, y, color.Gray{uint8(c.Rnd.IntN(int(c.Rule.States)))})
		}
	}
}

// Get the top-left neighbor value
func (c *CCA) getTopLeftNeighbor(x, y int) uint8 {
	w, h := c.Current.Rect.Dx(), c.Current.Rect.Dy()
	nx := (x - 1 + w) % w
	ny := (y - 1 + h) % h
	idx := c.Current.PixOffset(nx, ny)
	return c.Current.Pix[idx]
}

// Get the average value of von Neumann neighbors
func (c *CCA) getVonNeumannAverage(x, y int) uint8 {
	w, h := c.Current.Rect.Dx(), c.Current.Rect.Dy()

	sum := 0
	count := 0
	for _, offset := range NeumannNeighborsOffsets {
		nx := (x + offset.X + w) % w
		ny := (y + offset.Y + h) % h
		idx := c.Current.PixOffset(nx, ny)
		sum += int(c.Current.Pix[idx])
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
	w, h := c.Current.Rect.Dx(), c.Current.Rect.Dy()
	c.Neighbors = make([][]int, len(c.Current.Pix))
	for y := range h {
		for x := range w {
			index := c.Current.PixOffset(x, y)
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
					neighborIndex := c.Current.PixOffset(nx, ny)
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
	case LifeLike:
		c.lifeStep()
	}
}

// cyclicStep implements the original Cyclic rule
func (c *CCA) cyclicStep() {
	w, h := c.Current.Rect.Dx(), c.Current.Rect.Dy()
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
				idx := c.Current.PixOffset(x, y)
				currentState := c.Current.Pix[idx]
				targetState := (currentState + 1) % c.Rule.States
				count := 0

				for _, nIdx := range c.Neighbors[idx] {
					if c.Current.Pix[nIdx] == targetState {
						count++
						if count >= c.Rule.Threshold {
							break
						}
					}
				}

				if count >= c.Rule.Threshold {
					c.Next.Pix[idx] = targetState
				} else {
					c.Next.Pix[idx] = currentState
				}
			}
		}(y)
	}

	wg.Wait()

	// Swap grid and buffer
	c.Current, c.Next = c.Next, c.Current
}

// dunesStep implements the Dunes rule
func (c *CCA) dunesStep() {
	w, h := c.Current.Rect.Dx(), c.Current.Rect.Dy()
	var wg sync.WaitGroup
	wg.Add(h)

	// Update grid to the next generation concurrently, row by row
	for y := range h {
		// Capture current y
		go func(y int) {
			defer wg.Done()
			for x := range w {
				idx := c.Current.PixOffset(x, y)
				currentState := c.Current.Pix[idx]

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

				c.Next.Pix[idx] = newState
			}
		}(y)
	}

	wg.Wait()

	// Swap grid and buffer
	c.Current, c.Next = c.Next, c.Current
}

func (c *CCA) lifeProcessRow(y int, wg *sync.WaitGroup) {
	defer wg.Done()

	for x := range c.Current.Rect.Dx() {
		neighbors := CountMooreNeighbors(c.Current, x, y)
		current := c.Current.GrayAt(x, y).Y

		if current > 0 {
			if slices.Contains(c.Rule.Survive, neighbors) {
				if c.LifeLikeAgeEnabled {
					// Yaşayan hücrelerin rengi kademeli olarak azalır
					newColor := uint8(max(0, int(current)-(255/c.LifeLikeMaxAge)))
					c.Next.SetGray(x, y, color.Gray{Y: newColor})
				} else {
					c.Next.SetGray(x, y, color.Gray{Y: 255})
				}

			} else {
				c.Next.SetGray(x, y, color.Gray{Y: 0})
			}
		} else {
			if slices.Contains(c.Rule.Burn, neighbors) {
				// Yeni doğan hücreler tam parlaklıkta başlar
				c.Next.SetGray(x, y, color.Gray{Y: 255})
			} else {
				c.Next.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
}

// CountMooreNeighbors returns the number of alive neighbors of a cell in a Moore neighborhood
// Moore neighborhood includes the 8 surrounding cells
func CountMooreNeighbors(im *image.Gray, x, y int) int {
	count := 0
	for _, off := range MooreNeighborsOffsets {
		nx, ny := x+off.X, y+off.Y
		if nx >= 0 && nx < im.Bounds().Dx() && ny >= 0 && ny < im.Bounds().Dy() {
			if im.GrayAt(nx, ny).Y > 0 {
				count++
			}
		}
	}
	return count
}

func (c *CCA) lifeStep() {
	var wg sync.WaitGroup
	for y := range c.Current.Rect.Dy() {
		wg.Add(1)
		go c.lifeProcessRow(y, &wg)
	}
	wg.Wait()
	c.Current, c.Next = c.Next, c.Current
}

// RenderedImage returns a new image.Gray with color-mapped pixel values
func (c *CCA) RenderedImage() *image.Gray {
	if c.Rule.Type == LifeLike {
		return c.Current
	}
	ColorMap := make([]uint8, c.Rule.States)
	// Initialize the color map to evenly spread values across 0-255 range
	for i := range ColorMap {
		ColorMap[i] = uint8((i * 255) / (int(c.Rule.States) - 1))
	}

	w, h := c.Next.Rect.Dx(), c.Next.Rect.Dy()
	img := image.NewGray(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			idx := c.Next.PixOffset(x, y)
			state := c.Next.Pix[idx]
			img.Pix[idx] = ColorMap[state]
		}
	}
	return img
}

// RenderedImage returns a new image.Gray with color-mapped pixel values
func (c *CCA) GradientMap(g colorgrad.Gradient) image.Image {
	if c.Rule.Type == LifeLike {
		return c.Current
	}
	palet := g.Colors(uint(c.Rule.States))

	w, h := c.Next.Rect.Dx(), c.Next.Rect.Dy()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := range h {
		for x := range w {
			idx := c.Next.PixOffset(x, y)
			state := int(c.Next.Pix[idx])
			img.Set(x, y, palet[state])
		}
	}
	return img
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
func LifeRule(b, s []int) Rule {
	return Rule{
		Type:    LifeLike,
		Burn:    b,
		Survive: s,
		States:  2,
	}
}

func (c *CCA) RandomFillInRect(r image.Rectangle, full bool) {
	if full {
		r = c.Current.Bounds()
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if c.Rnd.IntN(2) == 0 {
				c.Current.SetGray(x, y, color.Gray{Y: 255})
			} else {
				c.Current.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

}

func (ca *CCA) FillWithRandom_0_255() {
	ca.RandomFillInRect(ca.Current.Bounds(), true)
}

func B(nums ...int) []int {
	return nums
}

func S(nums ...int) []int {
	return nums
}

func Birth(im *image.Gray, x, y int) {
	for _, off := range MooreNeighborsOffsets {
		im.SetGray(x+off.X, y+off.Y, color.Gray{Y: 255})
	}
}
