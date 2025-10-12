package ppu

const (
	fineYMask           uint16 = 0b111000000000000
	nametableMask       uint16 = 0b000110000000000
	nametableX          uint16 = 0b000010000000000
	nametableY          uint16 = 0b000100000000000
	coarseXMask         uint16 = 0b000000000011111
	coarseYMask         uint16 = 0b000001111100000
	addrNametableOffset uint16 = 0b010000000000000
	addrAttrtableOffset uint16 = 0b010001111000000
	fineYPosition       uint16 = 12
	nametablePosition   uint16 = 10
	coarseXPosition     uint16 = 0
	coarseYPosition     uint16 = 5
)

type vRegister struct {
	Value uint16
}

func (r *vRegister) CoarseX() uint16 {
	return (r.Value & coarseXMask) >> coarseXPosition
}

func (r *vRegister) IncrementX() {
	coarseX := r.CoarseX()
	if coarseX == 31 {
		coarseX = 0
		r.flipNametableX()
	} else {
		coarseX++
	}
	r.setCoarseX(coarseX)
}

func (r *vRegister) flipNametableX() {
	r.Value ^= nametableX
}

func (r *vRegister) setCoarseX(value uint16) {
	r.Value &= ^coarseXMask
	r.Value |= value & coarseXMask
}

func (r *vRegister) CoarseY() uint16 {
	return (r.Value & coarseYMask) >> coarseYPosition
}

func (r *vRegister) IncrementY() {
	fineY := r.FineY()
	if fineY < 7 {
		r.setFineY(fineY + 1)
		return
	}
	r.setFineY(0)
	coarseY := r.CoarseY()
	switch coarseY {
	case 29:
		r.setCoarseY(0)
		r.flipNametableY()
	case 31:
		r.setCoarseY(0)
	default:
		r.setCoarseY(coarseY + 1)
	}

}

func (r *vRegister) flipNametableY() {
	r.Value ^= nametableY
}

func (r *vRegister) setCoarseY(value uint16) {
	r.Value &= ^coarseYMask
	r.Value |= (value & (coarseYMask >> coarseYPosition)) << coarseYPosition
}

func (r *vRegister) FineY() uint16 {
	return (r.Value & fineYMask) >> fineYPosition
}

func (r *vRegister) setFineY(value uint16) {
	r.Value &= ^fineYMask
	r.Value |= (value & (fineYMask >> fineYPosition)) << fineYPosition
}

func (r *vRegister) Nametable() uint16 {
	return (r.Value & nametableMask) >> nametablePosition
}

func (r *vRegister) NametableAddress() uint16 {
	value := r.Value &^ fineYMask
	return addrNametableOffset | value
}

func (r *vRegister) AttrTableAddress() uint16 {
	coarseX := r.CoarseX() >> 2
	coarseY := r.CoarseY() >> 2
	nametable := r.Value & nametableMask
	value := addrAttrtableOffset | nametable | coarseY<<3 | coarseX
	return value
}

func (r *vRegister) AttrTableBytePart() uint16 {
	lo := (r.CoarseX() >> 1) & 0b01
	hi := r.CoarseY() & 0b10
	return hi | lo
}

func (r *vRegister) SetVerticalBits(value uint16) {
	mask := nametableY | fineYMask | coarseYMask
	verticalBits := value & mask
	r.Value &= ^mask
	r.Value |= verticalBits
}

func (r *vRegister) SetHorizontalBits(value uint16) {
	mask := nametableX | coarseXMask
	horizontalBits := value & mask
	r.Value &= ^mask
	r.Value |= horizontalBits
}
