package window

import (
	"image/color"
	"log"
	"unsafe"

	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	buttonAMask      uint8 = 0b00000001
	buttonBMask      uint8 = 0b00000010
	buttonSelectMask uint8 = 0b00000100
	buttonStartMask  uint8 = 0b00001000
	buttonUpMask     uint8 = 0b00010000
	buttonDownMask   uint8 = 0b00100000
	buttonLeftMask   uint8 = 0b01000000
	buttonRightMask  uint8 = 0b10000000
)

type WindowSize struct {
	Width  int
	Height int
}

func (s *WindowSize) Area() int {
	return s.Width * s.Height
}

type Window struct {
	closeChannel   chan struct{}
	joypadOne      *joypad.Joypad
	joypadOneState uint8
	joypadTwo      *joypad.Joypad
	joypadTwoState uint8
	window         *sdl.Window
	renderer       *sdl.Renderer
	texture        *sdl.Texture
	size           WindowSize
	imageChannel   chan []color.RGBA
}

func NewWindow(
	size WindowSize,
	joypadOne *joypad.Joypad,
	joypadTwo *joypad.Joypad,
	imageChannel chan []color.RGBA,
) *Window {
	closeChannel := make(chan struct{})
	return &Window{
		closeChannel: closeChannel,
		joypadOne:    joypadOne,
		joypadTwo:    joypadTwo,
		imageChannel: imageChannel,
		size:         size,
	}
}

func (w *Window) Width() int {
	return w.size.Width
}

func (w *Window) Height() int {
	return w.size.Height
}

func (w *Window) Start() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("failed to init sdl: %v", err)
	}
	defer sdl.Quit()

	w.window = w.createWindow()
	defer w.window.Destroy()

	w.renderer = w.createRenderer()
	defer w.renderer.Destroy()

	w.texture = w.createTexture()
	defer w.texture.Destroy()
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event := event.(type) {
			case *sdl.QuitEvent:
				w.closeChannel <- struct{}{}
				return
			case *sdl.KeyboardEvent:
				joypadOne, joypadTwo := w.readJoypadButton(event)
				w.joypadOneState |= joypadOne
				w.joypadTwoState |= joypadTwo
			}
		}

		w.joypadOneState &= w.joypadOne.Write(w.joypadOneState)
		w.joypadTwoState &= w.joypadTwo.Write(w.joypadTwoState)
		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_ESCAPE] != 0 {
			w.closeChannel <- struct{}{}
			return
		}

		select {
		case image := <-w.imageChannel:
			w.updateImage(image)
		default:
		}
		sdl.Delay(16)
	}
}

func (w *Window) createTexture() *sdl.Texture {
	texture, err := w.renderer.CreateTexture(
		sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(w.size.Width),
		int32(w.size.Height),
	)
	if err != nil {
		log.Fatalf("failed to create texture: %v", err)
	}
	return texture
}

func (w *Window) createRenderer() *sdl.Renderer {
	renderer, err := sdl.CreateRenderer(w.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("failed to create renderer: %v", err)
	}
	return renderer
}

func (w *Window) createWindow() *sdl.Window {
	window, err := sdl.CreateWindow(
		"NES",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(w.size.Width)*2,
		int32(w.size.Height)*2,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Fatalf("failed to create window: %v", err)
	}
	return window
}

func (w *Window) readJoypadButton(event *sdl.KeyboardEvent) (uint8, uint8) {
	if event.Type == sdl.KEYDOWN {
		switch event.Keysym.Sym {
		case sdl.K_SPACE:
			return buttonAMask, 0
		case sdl.K_LSHIFT:
			return buttonBMask, 0
		case sdl.K_BACKSPACE:
			return buttonSelectMask, buttonSelectMask
		case sdl.K_RETURN:
			return buttonStartMask, buttonStartMask
		case sdl.K_w:
			return buttonUpMask, 0
		case sdl.K_s:
			return buttonDownMask, 0
		case sdl.K_a:
			return buttonLeftMask, 0
		case sdl.K_d:
			return buttonRightMask, 0
		case sdl.K_1:
			return 0, buttonAMask
		case sdl.K_2:
			return 0, buttonBMask
		case sdl.K_UP:
			return 0, buttonUpMask
		case sdl.K_DOWN:
			return 0, buttonDownMask
		case sdl.K_LEFT:
			return 0, buttonLeftMask
		case sdl.K_RIGHT:
			return 0, buttonRightMask
		}
	}
	return 0, 0
}

func (w *Window) updateImage(colors []color.RGBA) {
	pixels := make([]byte, w.size.Area()*4)
	for i, color := range colors {
		offset := i * 4
		pixels[offset] = color.B
		pixels[offset+1] = color.G
		pixels[offset+2] = color.R
		pixels[offset+3] = color.A
	}

	w.texture.Update(nil, unsafe.Pointer(&pixels[0]), int(w.size.Width)*4)
	w.renderer.Clear()
	w.renderer.Copy(w.texture, nil, nil)
	w.renderer.Present()
}

func (w *Window) WaitUserExit() {
	<-w.closeChannel
}
