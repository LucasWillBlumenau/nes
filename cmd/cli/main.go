package main

import (
	"log"
	"os"

	"github.com/LucasWillBlumenau/nes/bus"
	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/LucasWillBlumenau/nes/ppu"
	"github.com/LucasWillBlumenau/nes/window"
)

func main() {
	path := readCliArgs()
	fp, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	cart, err := cartridge.LoadCartridgeFromReader(fp)
	if err != nil {
		log.Fatalln(err)
	}

	window := window.NewWindow(256, 240)
	ppu := ppu.NewPPU(window, cart.CharacterRom)
	bus := bus.NewBus(ppu, cart)
	cpu := cpu.NewCPU(bus)
	cpu.SetRomEntrypoint()

	go window.Start()
	for {
		cyclesTaken, err := cpu.Run()
		if err != nil {
			log.Fatalln(err)
		}
		image := ppu.ElapseCPUCycles(cyclesTaken)
		if image != nil {
			window.UpdateImageBuffer(image)
		}

	}
}

func readCliArgs() string {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("the program only supports a rom path as argument")
	}
	return args[0]
}
