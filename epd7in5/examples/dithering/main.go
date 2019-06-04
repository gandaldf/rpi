package main

import (
	"log"

	"github.com/MaxHalford/halfgone"
	"github.com/gandaldf/rpi/epd7in5"
)

func main() {
	log.Println("Starting...")
	epd, _ := epd7in5.New("P1_22", "P1_24", "P1_11", "P1_18")

	log.Println("Initializing the display...")
	epd.Init()

	log.Println("Clearing...")
	epd.Clear()

	// Test image
	log.Println("Opening nature.png test image...")
	var img, err = halfgone.LoadImage("../images/nature.png")
	if err != nil {
		log.Panic(err)
	}

	var gray = halfgone.ImageToGray(img)

	// Atkinson dithering
	var ad = halfgone.AtkinsonDitherer{}.Apply(gray)

	//halfgone.SaveImagePNG(ad, "dithered.png") // just for debug...

	log.Println("Displaying dithered image...")
	epd.Display(epd.Convert(ad))

	log.Println("Quitting...")
}
