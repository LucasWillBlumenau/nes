package ppu

type pixel struct {
	Palette uint8
	Color   uint8
}

type pixelsShiftRegister struct {
	buffer [16]pixel
	start  uint8
	end    uint8
}

func (p *pixelsShiftRegister) Buffer(color pixel) {
	index := p.end & 0b1111
	p.buffer[index] = color
	p.end++
}

func (p *pixelsShiftRegister) Unbuffer(offset uint8) pixel {
	currentIndex := (p.start + offset) & 0b1111
	p.start = (p.start + 1) & 0b1111
	return p.buffer[currentIndex]
}
