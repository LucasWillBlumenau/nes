package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/LucasWillBlumenau/nes/bus"
	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/LucasWillBlumenau/nes/ppu"
)

var pixelMap = [4]color.RGBA{
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // black
	{R: 0x55, G: 0x55, B: 0xFF, A: 0xFF}, // blue
	{R: 0xFF, G: 0x55, B: 0x55, A: 0xFF}, // red
	{R: 0xFF, G: 0xFF, B: 0xAA, A: 0xFF}, // yellow/cream
}

func main() {

	fp, err := os.Open("rom.nes")
	if err != nil {
		log.Fatalln(err)
	}

	cart, err := cartridge.LoadCartridgeFromReader(fp)
	if err != nil {
		log.Fatalln(err)
	}

	ppu := ppu.NewPPU(cart.CharacterRom)
	bus := bus.NewBus(ppu, cart)
	cpu := cpu.NewCPU(bus)

	go cpu.Run()
	go ppu.Run()

	time.Sleep(time.Hour)
}

func savePatternTable(cart *cartridge.Cartridge) {

	fmt.Printf("Length: %04x\n", len(cart.CharacterRom))

	type Tile [8][8]uint8

	var tiles []Tile
	offset := 0
	length := 8
	for offset < len(cart.CharacterRom) {
		leastSignificantBits := cart.CharacterRom[offset : offset+length]
		offset += length
		mostSignificantBits := cart.CharacterRom[offset : offset+length]
		offset += length

		var tile Tile
		for i := range length {
			loBits := leastSignificantBits[i]
			hiBits := mostSignificantBits[i]

			for j := range 8 {
				shiftSize := 7 - uint8(j)
				lo := (loBits >> shiftSize) & 0b00000001
				hi := (hiBits >> shiftSize) & 0b00000001
				tile[i][j] = (hi << 1) + lo
			}
		}
		tiles = append(tiles, tile)
	}

	image := image.NewRGBA(image.Rect(0, 0, 256, 128))

	for i, tile := range tiles {
		tileX := (i % 32) * 8
		tileY := (i / 32) * 8

		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				index := tile[y][x]
				image.Set(tileX+x, tileY+y, pixelMap[index])
			}
		}
	}

	out, _ := os.Create("output.png")
	defer out.Close()

	png.Encode(out, image)

	fmt.Printf("Found %d tiles\n", len(tiles))

	os.Exit(0)
}
