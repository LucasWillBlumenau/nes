package ppu

const (
	flipVerticallyMask   uint8 = 0b10000000
	flipHorizontallyMask uint8 = 0b01000000
	priorityMask         uint8 = 0b00100000
	paletteMask          uint8 = 0b00000011
)

type spriteAttributes struct {
	FlipVertically   bool
	FlipHorizontally bool
	HasPriority      bool
	Palette          uint8
}

func newSpriteAttributesFromByte(value uint8) spriteAttributes {
	flipVertically := (value & flipVerticallyMask) > 0
	flipHorizontally := (value & flipHorizontallyMask) > 0
	hasPriority := (value & priorityMask) == 0
	palette := (value & paletteMask)

	return spriteAttributes{
		FlipVertically:   flipVertically,
		FlipHorizontally: flipHorizontally,
		HasPriority:      hasPriority,
		Palette:          palette,
	}
}
