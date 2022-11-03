![7 5inch-e-paper-hat-4](https://user-images.githubusercontent.com/3932259/58586467-659e0380-825b-11e9-9942-f75c6dd7584f.jpg)

[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/gandaldf/rpi/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/gandaldf/rpi.svg?branch=master)](https://travis-ci.org/gandaldf/rpi)
[![Go Report Card](https://goreportcard.com/badge/github.com/gandaldf/rpi)](https://goreportcard.com/report/github.com/gandaldf/rpi)
[![Go Reference](https://pkg.go.dev/badge/github.com/gandaldf/rpi.svg)](https://pkg.go.dev/github.com/gandaldf/rpi)
[![Version](https://img.shields.io/github/tag/gandaldf/rpi.svg?color=blue&label=version)](https://github.com/gandaldf/rpi/releases)

# 7.5inch e-Paper
This is an interface for the Waveshare 7.5inch e-paper display ([wiki](https://www.waveshare.com/wiki/7.5inch_e-Paper_HAT)).

The GPIO and SPI communication is handled by the awesome **[Periph.io](https://periph.io/)** package; no CGO or other dependecy needed.

Tested on Raspberry Pi 3B / 3B+ / 4B with Raspbian Stretch.

For more information please check the _examples_ and _doc_ folders.

## Installation:
```
go get -u github.com/gandaldf/rpi/
```

## Load an image:
```golang
func main() {
	log.Println("Starting...")
	epd, _ := epd7in5.New("P1_22", "P1_24", "P1_11", "P1_18")

	log.Println("Initializing the display...")
	epd.Init()

	log.Println("Clearing...")
	epd.Clear()

	// Test image
	log.Println("Opening test image...")
	imgFile, err := os.Open("mypic.png")
	if err != nil {
		log.Panic(err)
	}

	defer imgFile.Close()

	img, err := png.Decode(imgFile)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Displaying test image...")
	epd.Display(epd.Convert(img))
}
```
#
For more information visit:

website: https://www.waveshare.com/

github: https://github.com/waveshare


