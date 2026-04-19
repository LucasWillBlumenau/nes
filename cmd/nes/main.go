package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/LucasWillBlumenau/nes/debug"
	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/LucasWillBlumenau/nes/nes"
	"github.com/LucasWillBlumenau/nes/window"
)

const (
	width  = 256
	height = 240
)

func main() {
	frames := make(chan image.RGBA)

	commands := make(chan debug.DebugOption)
	cartPath := readCliArgs()
	joypadOne := joypad.New()
	joypadTwo := joypad.New()
	scaleFactor := 2
	nes, err := nes.NewNES(
		frames,
		cartPath,
		scaleFactor,
		joypadOne,
		joypadTwo,
	)
	if err != nil {
		panic(err)
	}

	fmt.Print(`Debug options:

	1 - pause/reset frame generation
	2 - run next instruction
	3 - show/hide instructions
	4 - dump nametable
`)

	scaledWidth := width * scaleFactor
	scaledHeight := height * scaleFactor

	gameWindow := window.NewWindow(
		scaledWidth,
		scaledHeight,
		joypadOne,
		joypadTwo,
		frames,
		commands,
	)

	go nes.Run()
	go handleDebugOptions(nes, commands)
	gameWindow.Show()
}

func readCliArgs() string {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("the program only supports a rom path as argument")
	}
	return args[0]
}

func handleDebugOptions(nes *nes.NES, commands chan debug.DebugOption) {
	for cmd := range commands {
		switch cmd {
		case debug.DebugOptionPauseFrameGeneration:
			nes.Paused = !nes.Paused
		case debug.DebugOptionExecuteNextInstruction:
			nes.ExecuteNextInstruction()
		case debug.DebugOptionShowInstructions:
			nes.ShowInstructions = !nes.ShowInstructions
		case debug.DebugOptionDumpNametables:
			image := nes.DumpPPUNametables()
			saveImage(image)
		}
	}
}

func saveImage(image image.Image) {
	fileName := "nametable.png"
	out, _ := os.Create(fileName)
	defer out.Close()

	png.Encode(out, image)
}
