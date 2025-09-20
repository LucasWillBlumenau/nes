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
	f, err := os.Create("output.log")
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(f)

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
	ppu := ppu.NewPPU(cart.CharacterRom)
	bus := bus.NewBus(ppu, cart)
	cpu := cpu.NewCPU(bus)
	cpu.SetRomEntrypoint()

	go window.Start()
	go func() {
		for {
			cyclesTaken, err := cpu.Run()
			if err != nil {
				log.Fatalln(err)
			}
			for range cyclesTaken {
				ppu.RunStep()
			}
			if ppu.FrameDone() {
				image := ppu.GetGeneratedImage()
				window.UpdateImageBuffer(image)
			}

		}
	}()
	<-window.CloseChannel
}

func readCliArgs() string {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("the program only supports a rom path as argument")
	}
	return args[0]
}
