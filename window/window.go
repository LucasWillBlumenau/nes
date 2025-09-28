package window

import (
	"fmt"
	"image/color"
	"log"
	"unsafe"

	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	buttonAMask      uint8 = 0b11111110
	buttonBMask      uint8 = 0b11111101
	buttonSelectMask uint8 = 0b11111011
	buttonStartMask  uint8 = 0b11110111
	buttonUpMask     uint8 = 0b11101111
	buttonDownMask   uint8 = 0b11011111
	buttonLeftMask   uint8 = 0b10111111
	buttonRightMask  uint8 = 0b01111111
)

type WindowSize struct {
	Width  int
	Height int
}

func (s *WindowSize) Area() int {
	return s.Width * s.Height
}

type Window struct {
	CloseChannel chan struct{}
	imageBuffer  chan []color.RGBA
	joypadOne    *joypad.Joypad
	joypadTwo    *joypad.Joypad
	window       *sdl.Window
	renderer     *sdl.Renderer
	texture      *sdl.Texture
	size         WindowSize
}

func NewWindow(
	size WindowSize,
	joypadOne *joypad.Joypad,
	joypadTwo *joypad.Joypad,
) *Window {
	imageBuffer := make(chan []color.RGBA, 120)
	closeChannel := make(chan struct{})
	return &Window{
		imageBuffer:  imageBuffer,
		CloseChannel: closeChannel,
		joypadOne:    joypadOne,
		joypadTwo:    joypadTwo,
		size:         size,
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
		var joypadOneValue uint8 = buttonSelectMask // 0xFF
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event := event.(type) {
			case *sdl.QuitEvent:
				w.CloseChannel <- struct{}{}
				return
			case *sdl.KeyboardEvent:
				joypadOneValue &= w.readJoypadButton(event)
			}
		}
		w.joypadOne.Write(joypadOneValue)
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

func (w *Window) readJoypadButton(event *sdl.KeyboardEvent) uint8 {
	if event.Type == sdl.KEYDOWN {
		switch event.Keysym.Sym {
		case sdl.K_1:
			fmt.Println("Button A clicked")
			return buttonAMask
		case sdl.K_2:
			fmt.Println("Button B clicked")
			return buttonBMask
		case sdl.K_BACKSPACE:
			fmt.Println("Button Select clicked")
			return buttonSelectMask
		case sdl.K_RETURN:
			fmt.Println("Button Start clicked")
			return buttonStartMask
		case sdl.K_w:
			fmt.Println("ButtonU Up clicked")
			return buttonUpMask
		case sdl.K_s:
			fmt.Println("Button Down clicked")
			return buttonDownMask
		case sdl.K_a:
			fmt.Println("Button Left clicked")
			return buttonLeftMask
		case sdl.K_d:
			fmt.Println("Button Right clicked")
			return buttonRightMask
		}
	}
	return 0xFF
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
