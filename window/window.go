package window

import (
	"image/color"
	"log"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Window struct {
	CloseChannel chan bool
	imageBuffer  chan []color.RGBA
	window       *sdl.Window
	renderer     *sdl.Renderer
	texture      *sdl.Texture
	width        int
	height       int
}

func NewWindow(width int, height int) *Window {
	imageBuffer := make(chan []color.RGBA, 10)
	closeChannel := make(chan bool)
	return &Window{
		width:        width,
		height:       height,
		imageBuffer:  imageBuffer,
		CloseChannel: closeChannel,
	}
}

func (w *Window) Width() int {
	return w.width
}

func (w *Window) Height() int {
	return w.height
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
			switch event.(type) {
			case *sdl.QuitEvent:
				w.CloseChannel <- true
				return
			}
		}

		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_ESCAPE] != 0 {
			w.CloseChannel <- true
			return
		}

		select {
		case image := <-w.imageBuffer:
			w.updateImage(image)
		default:
		}
		sdl.Delay(16)
	}

}

func (w *Window) updateImage(colors []color.RGBA) {
	pixels := make([]byte, w.height*w.width*4)
	for i, color := range colors {
		offset := i * 4
		pixels[offset] = color.B
		pixels[offset+1] = color.G
		pixels[offset+2] = color.R
		pixels[offset+3] = color.A
	}

	w.texture.Update(nil, unsafe.Pointer(&pixels[0]), int(w.width)*4)
	w.renderer.Clear()
	w.renderer.Copy(w.texture, nil, nil)
	w.renderer.Present()
}

func (w *Window) createTexture() *sdl.Texture {
	texture, err := w.renderer.CreateTexture(
		sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(w.width),
		int32(w.height),
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
		int32(w.width)*2,
		int32(w.height)*2,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Fatalf("failed to create window: %v", err)
	}
	return window
}
