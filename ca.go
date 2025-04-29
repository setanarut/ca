package ca

import (
	"image"
	"image/color"
	"math/rand/v2"
	"sync"

	"slices"
)

var MooreNeighborsOffsets = [8]image.Point{
	{-1, -1}, {0, -1}, {1, -1}, {-1, 0},
	{1, 0}, {-1, 1}, {0, 1}, {1, 1}}

var NeumannNeighborsOffsets = [4]image.Point{
	{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

type CA struct {
	W, H          int
	Current, Next *image.Gray
	Burn          []int
	Survive       []int
}

// NewCA kurallar
//
// burn: Ölü bir hücre bu sayıda komşulara sahipse canlanır
//
// survive: Canlı bir hücre bu sayıda komşulara sahipse hayatta kalır
func NewCA(width, height int, burn, survive []int) *CA {
	current := image.NewGray(image.Rect(0, 0, width, height))
	next := image.NewGray(image.Rect(0, 0, width, height))
	return &CA{
		W:       width,
		H:       height,
		Current: current,
		Next:    next,
		Burn:    burn,
		Survive: survive,
	}
}

// CountMooreNeighbors returns the number of alive neighbors of a cell in a Moore neighborhood
// Moore neighborhood includes the 8 surrounding cells
func CountMooreNeighbors(im *image.Gray, x, y int) int {
	count := 0
	for _, off := range MooreNeighborsOffsets {
		nx, ny := x+off.X, y+off.Y
		if nx >= 0 && nx < im.Bounds().Dx() && ny >= 0 && ny < im.Bounds().Dy() {
			if im.GrayAt(nx, ny).Y == 255 {
				count++
			}
		}
	}
	return count
}

// ...existing code...

func (ca *CA) processRow(y int, wg *sync.WaitGroup) {
	defer wg.Done()

	for x := range ca.W {
		neighbors := CountMooreNeighbors(ca.Current, x, y)
		current := ca.Current.GrayAt(x, y).Y

		if current == 255 {
			survives := slices.Contains(ca.Survive, neighbors)
			if survives {
				ca.Next.SetGray(x, y, color.Gray{Y: 255})
			} else {
				ca.Next.SetGray(x, y, color.Gray{Y: 0})
			}
		} else {
			burns := slices.Contains(ca.Burn, neighbors)
			if burns {
				ca.Next.SetGray(x, y, color.Gray{Y: 255})
			} else {
				ca.Next.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
}

func (ca *CA) Step() {
	var wg sync.WaitGroup

	for y := range ca.H {
		wg.Add(1)
		go ca.processRow(y, &wg)
	}

	wg.Wait()
	ca.Current, ca.Next = ca.Next, ca.Current
}

func RandomFillInRect(im *image.Gray, r image.Rectangle, full bool) {
	if full {
		r = im.Bounds()
	}
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			if rand.IntN(2) == 0 {
				im.SetGray(x, y, color.Gray{Y: 255})
			} else {
				im.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

}
