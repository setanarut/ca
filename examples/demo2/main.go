package main

import (
	"ca"
	"fmt"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/mazznoer/colorgrad"
)

const (
	outputDir  = "frames"
	frameCount = (30 * 60) * 2 // dakika
	frameRate  = 30
	w, h       = 640, 480
	states     = 24
	// w, h = 854, 480
)

func main() {
	// Çıktı klasörünü oluştur
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Output klasörü oluşturulamadı: %v", err)
	}

	otomat := ca.NewAutomaton(w, h, 3, ca.DunesRule(states))
	otomat.FillWithRandomStates()

	// gradient := colorgrad.YlOrBr()

	// grad, _ := colorgrad.NewGradient().
	// 	Colors(
	// 		rgb(150, 92, 0),
	// 		rgb(205, 140, 0),
	// 		rgb(255, 209, 135),
	// 		rgb(205, 140, 0),
	// 		rgb(150, 92, 0),
	// 	).
	// 	Build()

	// grad, _ := colorgrad.NewGradient().
	// 	Colors(
	// 		rgb(255, 209, 135),
	// 		rgb(205, 140, 0),
	// 		rgb(150, 92, 0),
	// 	).
	// 	Build()

	for frame := range frameCount {

		for range 3 {
			otomat.Step()
		}

		// img := otomat.GradientMap(grad)
		img := otomat.RenderedImage()

		fileName := filepath.Join(outputDir, fmt.Sprintf("%04d.png", frame))
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("Dosya oluşturma hatası %s: %v", fileName, err)
		}

		if err := png.Encode(file, img); err != nil {
			log.Fatalf("PNG encode hatası: %v", err)
		}

		file.Close()
		percent := float64(frame+1) / float64(frameCount) * 100
		fmt.Printf("\rFrame: %d/%d (%.1f%%)", frame+1, frameCount, percent)
	}

	log.Printf("\nKareler '%s' klasörüne kaydedildi.", outputDir)
	imgio.Save("sonKare"+strconv.Itoa(states)+"_state.png", otomat.Next, imgio.PNGEncoder())
}

func rgb(r, g, b uint8) colorgrad.Color {
	return colorgrad.Rgb8(r, g, b, 255)
}
