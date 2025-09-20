package ppu

const (
	controlBaseNameTableAddr          = 0b00000011
	controlIncrementSize              = 0b00000100 // 0 = 1, 1 = 32
	controlSpritePatternTableAddr     = 0b00001000
	controlBackgroundPatternTableAddr = 0b00010000 // 0 = 0x0000, 1 = 0x1000
	controlSpriteSize                 = 0b00100000 // 0 = 8x8, 1 = 8x16
	controlMasterSlaveSelect          = 0b01000000
	controlPortNmiEnabled             = 0b10000000
)

type ppuControl struct {
	nametable                  uint16
	nametableOffset            uint16
	incrementSize              uint16
	spritePatternTableAddr     uint16
	backgroundPatternTableAddr uint16
	spriteSizeSet              bool
	masterSlave                bool
	nmiEnabled                 bool
}

func newPPUControl(value uint8) ppuControl {
	nametable := uint16(value & controlBaseNameTableAddr)
	var incrementSize uint16 = 1
	if (value & controlIncrementSize) > 0 {
		incrementSize = 32
	}
	var spritePatternTableAddr uint16 = 0x0000
	if (value & controlSpritePatternTableAddr) > 1 {
		spritePatternTableAddr = 0x1000
	}
	var backgroundPatternTableAddr uint16 = 0x0000
	if (value & controlBackgroundPatternTableAddr) > 1 {
		backgroundPatternTableAddr = 0x1000
	}
	spriteSizeSet := (value & controlSpriteSize) > 0
	masterSlave := (value & controlMasterSlaveSelect) > 0
	nmiEnabled := (value & controlPortNmiEnabled) > 0
	return ppuControl{
		nametable:                  nametable,
		nametableOffset:            (0b1000 | nametable) << 10,
		incrementSize:              incrementSize,
		spritePatternTableAddr:     spritePatternTableAddr,
		backgroundPatternTableAddr: backgroundPatternTableAddr,
		spriteSizeSet:              spriteSizeSet,
		masterSlave:                masterSlave,
		nmiEnabled:                 nmiEnabled,
	}
}
