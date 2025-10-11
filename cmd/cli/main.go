package main

import (
	"image"
	"log"
	"os"
	"runtime"

	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/LucasWillBlumenau/nes/nes"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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

func init() {
	runtime.LockOSThread()
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	var joypadOneState uint8 = 0
	var joypadTwoState uint8 = 0

	window, texture := setupWindow()
	frames := make(chan image.RGBA)

	cartPath := readCliArgs()
	joypadOne := joypad.New()
	joypadTwo := joypad.New()
	nes, err := nes.NewNES(
		frames,
		cartPath,
		2,
		joypadOne,
		joypadTwo,
	)
	if err != nil {
		panic(err)
	}

	go nes.Run()
	for !window.ShouldClose() {
		var width, heigth = window.GetSize()
		select {
		case img := <-nes.Frames:

			width32 := int32(width)
			heigth32 := int32(heigth)
			gl.BindTexture(gl.TEXTURE_2D, texture)
			gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, width32, heigth32, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
			gl.BlitFramebuffer(0, 0, width32, heigth32, 0, 0, width32, heigth32, gl.COLOR_BUFFER_BIT, gl.LINEAR)
			window.SwapBuffers()
			glfw.PollEvents()
		default:
			joypadOneState, joypadTwoState = readJoypadButton(window)
			joypadOneState &= joypadOne.Write(joypadOneState)
			joypadTwoState &= joypadTwo.Write(joypadTwoState)
		}
	}
}

func setupWindow() (*glfw.Window, uint32) {
	window, err := glfw.CreateWindow(512, 480, "NES", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	err = gl.Init()
	if err != nil {
		panic(err)
	}

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.BindImageTexture(0, texture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA8)
	var framebuffer uint32
	gl.GenFramebuffers(1, &framebuffer)
	gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, framebuffer)
	gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	return window, texture
}

func readJoypadButton(window *glfw.Window) (uint8, uint8) {
	var joypadOneState uint8 = 0
	var joypadTwoState uint8 = 0

	masks := []uint8{
		buttonAMask,
		buttonBMask,
		buttonUpMask,
		buttonDownMask,
		buttonLeftMask,
		buttonRightMask,
	}
	joypadOneKeys := []glfw.Key{
		glfw.KeySpace,
		glfw.KeyLeftShift,
		glfw.KeyW,
		glfw.KeyS,
		glfw.KeyA,
		glfw.KeyD,
	}

	joypadTwoKeys := []glfw.Key{
		glfw.Key1,
		glfw.Key2,
		glfw.KeyUp,
		glfw.KeyDown,
		glfw.KeyLeft,
		glfw.KeyRight,
	}

	for i := range masks {
		mask := masks[i]
		if window.GetKey(joypadOneKeys[i]) == glfw.Press {
			joypadOneState |= mask
		}
		if window.GetKey(joypadTwoKeys[i]) == glfw.Press {
			joypadTwoState |= mask
		}
	}

	if window.GetKey(glfw.KeyEnter) == glfw.Press {
		joypadOneState |= buttonStartMask
	}
	if window.GetKey(glfw.KeyBackspace) == glfw.Press {
		joypadOneState |= buttonSelectMask
	}
	return joypadOneState, joypadTwoState
}

func readCliArgs() string {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("the program only supports a rom path as argument")
	}
	return args[0]
}
