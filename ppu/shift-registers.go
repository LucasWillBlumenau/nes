package ppu

import "image/color"

type pixelsShiftRegister struct {
	buffer [16]color.RGBA
	start  uint8
	end    uint8
}

func (p *pixelsShiftRegister) Buffer(color color.RGBA) {
	index := p.end & 0b1111
	p.buffer[index] = color
	p.end++
}

func (p *pixelsShiftRegister) Unbuffer(offset uint8) color.RGBA {
	currentIndex := (p.start + offset) & 0b1111
	p.start = (p.start + 1) & 0b1111
	return p.buffer[currentIndex]
}
