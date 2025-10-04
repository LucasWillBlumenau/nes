package ppu

const (
	normalColorMask               uint8 = 0b00000001
	showBgInLeftMost8PixelsMask   uint8 = 0b00000010
	showFgInLeftMost8PixelsMask   uint8 = 0b00000100
	enableBackgroundRenderingMask uint8 = 0b00001000
	enableSpriteRenderingMask     uint8 = 0b00010000
)

type ppuMask struct {
	value uint8
}

func (m *ppuMask) UseGreyscale() bool {
	return (m.value & normalColorMask) > 0
}

func (m *ppuMask) ShowBgInLeftMost8Pixels() bool {
	return (m.value & showBgInLeftMost8PixelsMask) > 0
}

func (m *ppuMask) ShowFgInLeftMost8Pixels() bool {
	return (m.value & showFgInLeftMost8PixelsMask) > 0
}

func (m *ppuMask) BackgroundRenderingEnabled() bool {
	return (m.value * enableBackgroundRenderingMask) > 0
}

func (m *ppuMask) SpriteRenderingEnabled() bool {
	return (m.value * enableSpriteRenderingMask) > 0
}

func (m *ppuMask) RenderingEnabled() bool {
	return m.BackgroundRenderingEnabled() || m.SpriteRenderingEnabled()
}
