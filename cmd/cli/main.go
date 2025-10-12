package main

import (
	"fmt"
	"image"
	"log"
	"os"

	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/LucasWillBlumenau/nes/nes"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	buttonAIndex      = 0
	buttonBIndex      = 1
	buttonSelectIndex = 2
	buttonStartIndex  = 3
	buttonUpIndex     = 4
	buttonDownIndex   = 5
	buttonLeftIndex   = 6
	buttonRightIndex  = 7
)

var joypadButtonsMap = map[uint8]uint8{
	0: buttonAIndex,
	1: buttonBIndex,
	2: buttonBIndex,
	3: buttonAIndex,
	6: buttonSelectIndex,
	7: buttonStartIndex,
}

var joypadDPadMap = map[uint8][]uint8{
	sdl.HAT_UP:        {buttonUpIndex},
	sdl.HAT_LEFTUP:    {buttonUpIndex, buttonLeftIndex},
	sdl.HAT_RIGHTUP:   {buttonUpIndex, buttonRightIndex},
	sdl.HAT_RIGHT:     {buttonRightIndex},
	sdl.HAT_DOWN:      {buttonDownIndex},
	sdl.HAT_LEFTDOWN:  {buttonDownIndex},
	sdl.HAT_RIGHTDOWN: {buttonDownIndex},
	sdl.HAT_LEFT:      {buttonLeftIndex},
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("error initializing SDL: %s", err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("NES", 0, 0, 512, 480, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("error creating window: %s", err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		log.Fatalf("error creating window: %s", err)
	}
	defer surface.Free()

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

	for i := range sdl.NumJoysticks() {
		if sdl.IsGameController(i) {
			controller := sdl.GameControllerOpen(i)
			fmt.Printf("controller %s added\n", controller.Name())
		}
	}

	go nes.Run()
	for {
		select {
		case img := <-nes.Frames:
			var width, heigth = window.GetSize()
			for x := range int(width) {
				for y := range int(heigth) {
					surface.Set(x, y, img.At(x, y))
				}
			}
			window.UpdateSurface()
		default:
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.GetType() {
				case sdl.QUIT:
					return
				case sdl.KEYDOWN,
					sdl.KEYUP,
					sdl.JOYBUTTONDOWN,
					sdl.JOYBUTTONUP,
					sdl.JOYAXISMOTION,
					sdl.JOYHATMOTION:
					readJoypad(joypadOne, event)
				}
			}
		}
	}
}

func readJoypad(joypad *joypad.Joypad, event sdl.Event) {
	switch event := event.(type) {
	case *sdl.JoyButtonEvent:
		button, ok := joypadButtonsMap[event.Button]
		if !ok {
			return
		}
		if event.State == sdl.RELEASED {
			joypad.SetControl(button, false)
		} else {
			joypad.SetControl(button, true)
		}
	case *sdl.JoyHatEvent:
		var buttons []uint8
		var value bool
		if selectedButtons, ok := joypadDPadMap[event.Value]; ok {
			buttons = selectedButtons
			value = true
		} else {
			buttons = []uint8{buttonUpIndex, buttonDownIndex, buttonLeftIndex, buttonRightIndex}
			value = false
		}

		for _, button := range buttons {
			joypad.SetControl(button, value)
		}
	}
}

func readCliArgs() string {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("the program only supports a rom path as argument")
	}
	return args[0]
}
