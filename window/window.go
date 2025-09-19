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
	imageBuffer := make(chan []color.RGBA, 120)
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

	// Define the border color (you can change this to whatever color you want)
	borderColor := color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red border

	// Loop through each 8x8 tile
	for tileY := 0; tileY < w.height/8; tileY++ {
		for tileX := 0; tileX < w.width/8; tileX++ {
			// Calculate the starting pixel for the current tile
			tileStartX := tileX * 8
			tileStartY := tileY * 8

			// Loop through each pixel in the current 8x8 tile
			for y := 0; y < 8; y++ {
				for x := 0; x < 8; x++ {
					// Calculate the index of the pixel in the flat array
					index := (tileStartY+y)*w.width + (tileStartX + x)
					pixelIndex := index * 4

					// Check if this pixel is on the border (top, bottom, left, or right edge of the tile)
					if y == 0 || y == 7 || x == 0 || x == 7 {
						// Set border color
						pixels[pixelIndex] = borderColor.B
						pixels[pixelIndex+1] = borderColor.G
						pixels[pixelIndex+2] = borderColor.R
						pixels[pixelIndex+3] = borderColor.A
					} else {
						// Set the color from the original colors array
						color := colors[index]
						pixels[pixelIndex] = color.B
						pixels[pixelIndex+1] = color.G
						pixels[pixelIndex+2] = color.R
						pixels[pixelIndex+3] = color.A
					}
				}
			}
		}
	}

	// Update the texture with the pixel data
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
