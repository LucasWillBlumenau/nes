package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/window"
)

var pixelMap = [4]color.RGBA{
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // black
	{R: 0x55, G: 0x55, B: 0xFF, A: 0xFF}, // blue
	{R: 0xFF, G: 0x55, B: 0x55, A: 0xFF}, // red
	{R: 0xFF, G: 0xFF, B: 0xAA, A: 0xFF}, // yellow/cream
}

const (
	screenWidth  int = 256
	screenHeight int = 128
)

func main() {
	fp, err := os.Open("roms/rom.nes")
	if err != nil {
		log.Fatalln(err)
	}

	cart, err := cartridge.LoadCartridgeFromReader(fp)
	if err != nil {
		log.Fatalln(err)
	}

	window := window.NewWindow(screenWidth, screenHeight)

	tiles := makeTiles(cart.CharacterRom)
	colors := make([]color.RGBA, screenWidth*screenHeight)
	for i, tile := range tiles {
		tileX := (i % 32) * 8
		tileY := (i / 32) * 8

		for y := range 8 {
			for x := range 8 {
				colorIndex := tile[y][x]
				color := pixelMap[colorIndex]
				colors[(tileX+x)+(tileY+y)*screenWidth] = color
			}
		}
	}

	go window.Start()
	go func() {
		start := time.Now().UnixNano()
		frames := 0

		for {
			window.UpdateImageBuffer(colors)
			frames++
			now := time.Now().UnixNano()

			delta := now - start
			if delta >= int64(time.Second) {
				fmt.Printf("Frames generated: %d\n", frames)
				frames = 0
				start = time.Now().UnixNano()
			}

		}
	}()

	<-window.CloseChannel

}

type Tile [8][8]uint8

func makeTiles(buffer []uint8) []Tile {
	offset := 0
	length := 8

	tiles := make([]Tile, 0)

	for offset < len(buffer) {
		leastSignificantBytes := buffer[offset : offset+length]
		offset += length
		mostSignificantBytes := buffer[offset : offset+length]
		offset += length

		var tile Tile
		for i := range length {
			leastSignificantBits := leastSignificantBytes[i]
			mostSignificantBits := mostSignificantBytes[i]

			for j := range 8 {
				shiftSize := 7 - j
				bit0 := (leastSignificantBits >> shiftSize) & 0x01
				bit1 := (mostSignificantBits >> shiftSize) & 0x01
				tile[i][j] = uint8(bit1)<<1 | uint8(bit0)
			}

		}
		tiles = append(tiles, tile)
	}

	return tiles
}
