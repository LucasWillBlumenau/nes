package main

import (
	"image/color"
	"log"
	"os"

	"github.com/LucasWillBlumenau/nes/bus"
	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/LucasWillBlumenau/nes/ppu"
	"github.com/LucasWillBlumenau/nes/window"
)

func main() {

	f, err := os.Create("output.log")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	log.SetOutput(f)
	// os.Stdout = f
	romPath := readCliArgs()
	cart, err := cartridge.LoadCartridge(romPath)
	if err != nil {
		log.Fatalln(err)
	}

	joypadOne := joypad.New()
	joypadTwo := joypad.New()

	imageChannel := make(chan []color.RGBA, 1)
	windowSize := window.WindowSize{Width: 256, Height: 240}
	gameWindow := window.NewWindow(
		windowSize,
		joypadOne,
		joypadTwo,
		imageChannel,
	)

	ppuBus := ppu.NewPPUBus(cart)
	ppu := ppu.NewPPU(ppuBus, imageChannel)
	bus := bus.NewBus(ppu, cart, joypadOne, joypadTwo)
	cpu := cpu.NewCPU(bus)
	cpu.SetRomEntrypoint()
	// cpu.Pc = 0xC000

	go gameWindow.Start()
	go func() {
		for {
			cyclesTaken, err := cpu.Run()
			if err != nil {
				log.Fatalln(err)
			}
			ppuCycles := cyclesTaken * 3
			for range ppuCycles {
				ppu.RunStep()
			}

		}
	}()
	<-gameWindow.CloseChannel
	// time.Sleep(time.Second * 1000)
}

func readCliArgs() string {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("the program only supports a rom path as argument")
	}
	return args[0]
}
