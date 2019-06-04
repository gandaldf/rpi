package main

import (
	"log"

	"github.com/fogleman/gg"
	"github.com/gandaldf/rpi/epd7in5"
)

func main() {
	log.Println("Starting...")
	epd, _ := epd7in5.New("P1_22", "P1_24", "P1_11", "P1_18")

	log.Println("Initializing the display...")
	epd.Init()

	log.Println("Clearing...")
	epd.Clear()

	// Test image 1
	log.Println("Drawing...")
	img1, err := gg.LoadImage("../images/corners.png")
	if err != nil {
		log.Panic(err)
	}

	// Test image 2
	img2, err := gg.LoadImage("../images/100x100.png")
	if err != nil {
		log.Panic(err)
	}

	// New context
	cx := gg.NewContextForImage(img1)
	cx.Clear()

	// Place our test images
	cx.DrawImage(img1, 0, 0)
	cx.DrawImage(img2, 430, 46)

	// Draw a circle
	cx.SetRGB(0, 0, 0)
	cx.DrawCircle(80, 288, 50)
	cx.Stroke()

	// Draw a square
	cx.SetRGB(0, 0, 0)
	cx.DrawRectangle(180, 238, 100, 100)
	cx.Fill()

	// Print some text
	cx.SetRGB(0, 0, 0)
	myString := "Hello world!"
	sw, sh := cx.MeasureString(myString)
	cx.DrawString(myString, 480-(sw/2), 288-(sh/2))

	//cx.SavePNG("output.png") // just for debug...

	log.Println("Displaying output image...")
	epd.Display(epd.Convert(cx.Image()))

	log.Println("Quitting...")
}
