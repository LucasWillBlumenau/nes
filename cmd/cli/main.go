package main

import (
	"log"
	"os"
	"time"

	"github.com/LucasWillBlumenau/nes/bus"
	"github.com/LucasWillBlumenau/nes/cartridge"
	"github.com/LucasWillBlumenau/nes/cpu"
	"github.com/LucasWillBlumenau/nes/ppu"
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

	ppu := ppu.NewPPU(cart.CharacterRom)
	bus := bus.NewBus(ppu, cart)
	cpu := cpu.NewCPU(bus)

	// ppu.WritePPUAddrPort(0x20)
	// ppu.WritePPUAddrPort(0x00)
	// ppu.WritePPUDataPort(0x00)
	// ppu.WritePPUDataPort(0x01)
	// ppu.WritePPUDataPort(0x02)
	// ppu.WritePPUDataPort(0x03)
	// ppu.WritePPUDataPort(0x04)
	// ppu.WritePPUDataPort(0x05)
	// ppu.WritePPUDataPort(0x06)
	// ppu.WritePPUDataPort(0x07)

	// ppu.OutputCurrentNameTable()
	// return

	cpu.SetResetInterruptHandlerAddressAsEntrypoint()
	for ppu.VBlankCount < 60 {
		cyclesTaken, err := cpu.Run()
		if err != nil {
			log.Fatalln(err)
		}
		ppu.ElapseCPUCycles(cyclesTaken)
	}
	ppu.OutputCurrentNameTable()
}
