// Package epd7in5 is an interface for the Waveshare 7.5inch e-paper display (wiki).
//
// The GPIO and SPI communication is handled by the awesome Periph.io package; no CGO or other dependecy needed.
//
// Tested on Raspberry Pi 3B / 3B+ with Raspbian Stretch.
//
// For more information please check the examples and doc folders.
package epd7in5

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"time"

	"periph.io/x/periph/conn"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/host"
)

const (
	EPD_WIDTH  int = 640
	EPD_HEIGHT int = 384
)

const (
	PANEL_SETTING                  byte = 0x00
	POWER_SETTING                  byte = 0x01
	POWER_OFF                      byte = 0x02
	POWER_OFF_SEQUENCE_SETTING     byte = 0x03
	POWER_ON                       byte = 0x04
	POWER_ON_MEASURE               byte = 0x05
	BOOSTER_SOFT_START             byte = 0x06
	DEEP_SLEEP                     byte = 0x07
	DATA_START_TRANSMISSION_1      byte = 0x10
	DATA_STOP                      byte = 0x11
	DISPLAY_REFRESH                byte = 0x12
	IMAGE_PROCESS                  byte = 0x13
	LUT_FOR_VCOM                   byte = 0x20
	LUT_BLUE                       byte = 0x21
	LUT_WHITE                      byte = 0x22
	LUT_GRAY_1                     byte = 0x23
	LUT_GRAY_2                     byte = 0x24
	LUT_RED_0                      byte = 0x25
	LUT_RED_1                      byte = 0x26
	LUT_RED_2                      byte = 0x27
	LUT_RED_3                      byte = 0x28
	LUT_XON                        byte = 0x29
	PLL_CONTROL                    byte = 0x30
	TEMPERATURE_SENSOR_COMMAND     byte = 0x40
	TEMPERATURE_CALIBRATION        byte = 0x41
	TEMPERATURE_SENSOR_WRITE       byte = 0x42
	TEMPERATURE_SENSOR_READ        byte = 0x43
	VCOM_AND_DATA_INTERVAL_SETTING byte = 0x50
	LOW_POWER_DETECTION            byte = 0x51
	TCON_SETTING                   byte = 0x60
	TCON_RESOLUTION                byte = 0x61
	SPI_FLASH_CONTROL              byte = 0x65
	REVISION                       byte = 0x70
	GET_STATUS                     byte = 0x71
	AUTO_MEASUREMENT_VCOM          byte = 0x80
	READ_VCOM_VALUE                byte = 0x81
	VCM_DC_SETTING                 byte = 0x82
)

// Epd is a handle to the display controller.
type Epd struct {
	c          conn.Conn
	dc         gpio.PinOut
	cs         gpio.PinOut
	rst        gpio.PinOut
	busy       gpio.PinIO
	widthByte  int
	heightByte int
}

// New returns a Epd object that communicates over SPI to the display controller.
func New(dcPin, csPin, rstPin, busyPin string) (*Epd, error) {
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	// DC pin
	dc := gpioreg.ByName(dcPin)
	if dc == nil {
		return nil, errors.New("spi: failed to find DC pin")
	}

	if dc == gpio.INVALID {
		return nil, errors.New("epd: use nil for dc to use 3-wire mode, do not use gpio.INVALID")
	}

	if err := dc.Out(gpio.Low); err != nil {
		return nil, err
	}

	// CS pin
	cs := gpioreg.ByName(csPin)
	if cs == nil {
		return nil, errors.New("spi: failed to find CS pin")
	}

	if err := cs.Out(gpio.Low); err != nil {
		return nil, err
	}

	// RST pin
	rst := gpioreg.ByName(rstPin)
	if rst == nil {
		return nil, errors.New("spi: failed to find RST pin")
	}

	if err := rst.Out(gpio.Low); err != nil {
		return nil, err
	}

	// BUSY pin
	busy := gpioreg.ByName(busyPin)
	if busy == nil {
		return nil, errors.New("spi: failed to find BUSY pin")
	}

	if err := busy.In(gpio.PullDown, gpio.RisingEdge); err != nil {
		return nil, err
	}

	// SPI
	port, err := spireg.Open("")
	if err != nil {
		return nil, err
	}

	c, err := port.Connect(5*physic.MegaHertz, spi.Mode0, 8)
	if err != nil {
		port.Close()
		return nil, err
	}

	var widthByte, heightByte int

	if EPD_WIDTH%8 == 0 {
		widthByte = (EPD_WIDTH / 8)
	} else {
		widthByte = (EPD_WIDTH/8 + 1)
	}

	heightByte = EPD_HEIGHT

	e := &Epd{
		c:          c,
		dc:         dc,
		cs:         cs,
		rst:        rst,
		busy:       busy,
		widthByte:  widthByte,
		heightByte: heightByte,
	}

	return e, nil
}

// Reset can be also used to awaken the device.
func (e *Epd) Reset() {
	e.rst.Out(gpio.High)
	time.Sleep(200 * time.Millisecond)
	e.rst.Out(gpio.Low)
	time.Sleep(200 * time.Millisecond)
	e.rst.Out(gpio.High)
	time.Sleep(200 * time.Millisecond)
}

func (e *Epd) sendCommand(cmd byte) {
	e.dc.Out(gpio.Low)
	e.cs.Out(gpio.Low)
	e.c.Tx([]byte{cmd}, nil)
	e.cs.Out(gpio.High)
}

func (e *Epd) sendData(data byte) {
	e.dc.Out(gpio.High)
	e.cs.Out(gpio.Low)
	e.c.Tx([]byte{data}, nil)
	e.cs.Out(gpio.High)
}

func (e *Epd) waitUntilIdle() {
	for e.busy.Read() == gpio.Low {
		time.Sleep(100 * time.Millisecond)
	}
}

func (e *Epd) turnOnDisplay() {
	e.sendCommand(DISPLAY_REFRESH)
	time.Sleep(100 * time.Millisecond)
	e.waitUntilIdle()
}

// Init initializes the display config.
// It should be only used when you put the device to sleep and need to re-init the device.
func (e *Epd) Init() {
	e.Reset()

	e.sendCommand(POWER_SETTING)
	e.sendData(0x37)
	e.sendData(0x00)

	e.sendCommand(PANEL_SETTING)
	e.sendData(0xCF)
	e.sendData(0x08)

	e.sendCommand(BOOSTER_SOFT_START)
	e.sendData(0xc7)
	e.sendData(0xcc)
	e.sendData(0x28)

	e.sendCommand(POWER_ON)
	e.waitUntilIdle()

	e.sendCommand(PLL_CONTROL)
	e.sendData(0x3c)

	e.sendCommand(TEMPERATURE_CALIBRATION)
	e.sendData(0x00)

	e.sendCommand(VCOM_AND_DATA_INTERVAL_SETTING)
	e.sendData(0x77)

	e.sendCommand(TCON_SETTING)
	e.sendData(0x22)

	e.sendCommand(TCON_RESOLUTION)
	e.sendData(byte(EPD_WIDTH >> 8))
	e.sendData(byte(EPD_WIDTH & 0xff))
	e.sendData(byte(EPD_HEIGHT >> 8))
	e.sendData(byte(EPD_HEIGHT & 0xff))

	e.sendCommand(VCM_DC_SETTING)
	e.sendData(0x1E)

	e.sendCommand(0xe5)
	e.sendData(0x03)
}

// Clear clears the screen.
func (e *Epd) Clear() {
	e.sendCommand(DATA_START_TRANSMISSION_1)

	for j := 0; j < e.heightByte; j++ {
		for i := 0; i < e.widthByte; i++ {
			for k := 0; k < 4; k++ {
				e.sendData(0x33)
			}
		}
	}

	e.turnOnDisplay()
}

// Display takes a byte buffer and updates the screen.
func (e *Epd) Display(img []byte) {
	e.sendCommand(DATA_START_TRANSMISSION_1)

	for j := 0; j < e.heightByte; j++ {
		for i := 0; i < e.widthByte; i++ {
			dataBlack := ^img[i+j*e.widthByte]

			for k := 0; k < 8; k++ {
				var data byte

				if dataBlack&0x80 > 0 {
					data = 0x00
				} else {
					data = 0x03
				}

				data <<= 4
				dataBlack <<= 1
				k++

				if dataBlack&0x80 > 0 {
					data |= 0x00
				} else {
					data |= 0x03
				}

				dataBlack <<= 1

				e.sendData(data)
			}
		}
	}

	e.turnOnDisplay()
}

// Sleep put the display in power-saving mode.
// You can use Reset() to awaken and Init() to re-initialize the display.
func (e *Epd) Sleep() {
	e.sendCommand(POWER_OFF)
	e.waitUntilIdle()
	e.sendCommand(DEEP_SLEEP)
	e.sendData(0XA5)
}

// Convert converts the input image into a ready-to-display byte buffer.
func (e *Epd) Convert(img image.Image) []byte {
	var byteToSend byte = 0x00
	var bgColor = 1

	buffer := bytes.Repeat([]byte{0x00}, e.widthByte*e.heightByte)

	for j := 0; j < EPD_HEIGHT; j++ {
		for i := 0; i < EPD_WIDTH; i++ {
			bit := bgColor

			if i < img.Bounds().Dx() && j < img.Bounds().Dy() {
				bit = color.Palette([]color.Color{color.Black, color.White}).Index(img.At(i, j))
			}

			if bit == 1 {
				byteToSend |= 0x80 >> (uint32(i) % 8)
			}

			if i%8 == 7 {
				buffer[(i/8)+(j*e.widthByte)] = byteToSend
				byteToSend = 0x00
			}
		}
	}

	return buffer
}
