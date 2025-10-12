package main

import (
	"image"
	"log"
	"os"

	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/LucasWillBlumenau/nes/nes"
	"github.com/LucasWillBlumenau/nes/window"
)

const (
	width  = 512
	heigth = 240
)

func main() {

	frames := make(chan image.RGBA)
	cartPath := readCliArgs()
	joypadOne := joypad.New()
	joypadTwo := joypad.New()
	nes, err := nes.NewNES(
		frames,
		cartPath,
		2,
		joypadOne,
		joypadTwo,
	)
	if err != nil {
		panic(err)
	}

	window := window.NewWindow(
		width,
		heigth,
		joypadOne,
		joypadTwo,
		frames,
	)
	go nes.Run()
	window.Show()
}

func readCliArgs() string {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("the program only supports a rom path as argument")
	}
	return args[0]
}
