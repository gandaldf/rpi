package main

import (
	"image/png"
	"log"
	"os"
	"time"

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
	log.Println("Opening 7in5.png test image...")
	imgFile1, err := os.Open("../images/7in5.png")
	if err != nil {
		log.Panic(err)
	}

	defer imgFile1.Close()

	img1, err := png.Decode(imgFile1)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Displaying 7in5.png test image...")
	epd.Display(epd.Convert(img1))

	time.Sleep(5 * time.Second)

	// Test image 2
	log.Println("Opening 100x100.png test image...")
	imgFile2, err := os.Open("../images/100x100.png")
	if err != nil {
		log.Panic(err)
	}

	defer imgFile2.Close()

	img2, err := png.Decode(imgFile2)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Displaying 100x100.png test image...")
	epd.Display(epd.Convert(img2))

	time.Sleep(5 * time.Second)

	// Test image 3
	log.Println("Opening corners.png test image...")
	imgFile3, err := os.Open("../images/corners.png")
	if err != nil {
		log.Panic(err)
	}

	defer imgFile3.Close()

	img3, err := png.Decode(imgFile3)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Displaying corners.png test image...")
	epd.Display(epd.Convert(img3))

	log.Println("Quitting...")
}
