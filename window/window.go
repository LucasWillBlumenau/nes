package window

import (
	"fmt"
	"image/color"
	"log"
	"unsafe"

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
	CloseChannel     chan struct{}
	imageBuffer      chan []color.RGBA
	joypadOneChannel *uint8
	joypadTwoChannel *uint8
	window           *sdl.Window
	renderer         *sdl.Renderer
	texture          *sdl.Texture
	size             WindowSize
}

func NewWindow(
	size WindowSize,
	joypadOneChannel *uint8,
	joypadTwoChannel *uint8,
) *Window {
	imageBuffer := make(chan []color.RGBA, 120)
	closeChannel := make(chan struct{})
	return &Window{
		imageBuffer:      imageBuffer,
		CloseChannel:     closeChannel,
		joypadOneChannel: joypadOneChannel,
		joypadTwoChannel: joypadTwoChannel,
		size:             size,
	}
}

func (w *Window) Width() int {
	return w.size.Width
}

func (w *Window) Height() int {
	return w.size.Height
}

func (w *Window) UpdateImageBuffer(image []color.RGBA) {
	w.imageBuffer <- image
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
				w.CloseChannel <- struct{}{}
				return
			case *sdl.KeyboardEvent:
				w.writeToJoypad(event)
			}
		}

		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_ESCAPE] != 0 {
			w.CloseChannel <- struct{}{}
			return
		}

		select {
		case image := <-w.imageBuffer:
			w.updateImage(image)
			sdl.Delay(16)
		default:
		}
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

func (w *Window) writeToJoypad(event *sdl.KeyboardEvent) {
	switch event.Type {
	case sdl.KEYDOWN:
		switch event.Keysym.Sym {
		case sdl.K_1:
			*w.joypadOneChannel = buttonAMask
			fmt.Println("Button A clicked")
		case sdl.K_2:
			*w.joypadOneChannel = buttonBMask
			fmt.Println("Button B clicked")
		case sdl.K_BACKSPACE:
			*w.joypadOneChannel = buttonSelectMask
			fmt.Println("Button Select clicked")
		case sdl.K_KP_ENTER:
			*w.joypadOneChannel = buttonStartMask
			fmt.Println("Button Start clicked")
		case sdl.K_w:
			*w.joypadOneChannel = buttonUpMask
			fmt.Println("ButtonU Up clicked")
		case sdl.K_s:
			*w.joypadOneChannel = buttonDownMask
			fmt.Println("Button Down clicked")
		case sdl.K_a:
			*w.joypadOneChannel = buttonLeftMask
			fmt.Println("Button Left clicked")
		case sdl.K_d:
			*w.joypadOneChannel = buttonRightMask
			fmt.Println("Button Right clicked")
		default:
		}
	case sdl.KEYUP:
		*w.joypadOneChannel = 0
	}
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
