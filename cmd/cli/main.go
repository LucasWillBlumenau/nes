package main

import (
	"log"
	"os"
	"time"

	"github.com/LucasWillBlumenau/nes/bus"
	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/LucasWillBlumenau/nes/ppu"
	"github.com/LucasWillBlumenau/nes/window"
)

const clockDurationInNanoseconds time.Duration = 559

func main() {
	fp, err := os.Open("rom.nes")
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
	cpu.SetResetInterruptHandlerAddressAsEntrypoint()

	go window.Start()
	for {
		cyclesTaken, err := cpu.Run()
		if err != nil {
			log.Fatalln(err)
		}
		ppu.ElapseCPUCycles(cyclesTaken)
	}
}
