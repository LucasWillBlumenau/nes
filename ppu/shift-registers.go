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

const indexMask uint8 = 0b1111

func (p *pixelsShiftRegister) Buffer(color pixel) {
	index := p.end & indexMask
	p.buffer[index] = color
	p.end++
}

func (p *pixelsShiftRegister) Unbuffer(offset uint8) pixel {
	currentIndex := (p.start + offset) & indexMask
	p.start = (p.start + 1) & indexMask
	return p.buffer[currentIndex]
}
