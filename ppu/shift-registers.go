package ppu

type pixel struct {
	Palette uint8
	Color   uint8
}

type pixelsShiftRegister struct {
	buffer [272]pixel
	start  uint16
	end    uint16
	fineX  uint16
}

func (p *pixelsShiftRegister) SetFineX(fineX uint16) {
	p.fineX = fineX
}

func (p *pixelsShiftRegister) Reset() {
	p.end = 0
	p.start = 0
}

func (p *pixelsShiftRegister) Buffer(color pixel) {
	p.buffer[p.end] = color
	p.end++
}

func (p *pixelsShiftRegister) Unbuffer() pixel {
	value := p.buffer[p.start+p.fineX]
	p.start++
	return value
}
