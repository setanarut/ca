package main

import (
	"ca"
	"fmt"
	"image"
	"time"

	"github.com/anthonynsimon/bild/imgio"
)

func main() {

	// B3/S012345678	Life without Death
	// B3/23	Conway

	burn := []int{3}
	survive := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

	caa := ca.NewCA(800, 512, burn, survive)
	x0 := (caa.W - 10) / 2
	y0 := (caa.H - 10) / 2
	ca.RandomFillInRect(caa.Current, image.Rect(x0, y0, x0+10, y0+10), false)

	imgio.Save("output/ilkKare.png", caa.Current, imgio.PNGEncoder())

	now := time.Now()
	// 600 adım at
	for range 600 {
		caa.Step()
	}
	fmt.Println("Elapsed time:", time.Since(now).Seconds())

	imgio.Save("output/sonKare.png", caa.Current, imgio.PNGEncoder())

}
