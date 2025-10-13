package window

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/LucasWillBlumenau/nes/joypad"
	"github.com/veandco/go-sdl2/sdl"
)

var joypadButtonsMap = map[uint8]joypad.Button{
	0: joypad.ButtonA,
	1: joypad.ButtonB,
	2: joypad.ButtonB,
	3: joypad.ButtonA,
	6: joypad.ButtonSelect,
	7: joypad.ButtonStart,
}

var joypadDPadMap = map[uint8][]joypad.Button{
	sdl.HAT_UP:        {joypad.ButtonUp},
	sdl.HAT_LEFTUP:    {joypad.ButtonUp, joypad.ButtonLeft},
	sdl.HAT_RIGHTUP:   {joypad.ButtonUp, joypad.ButtonRight},
	sdl.HAT_RIGHT:     {joypad.ButtonRight},
	sdl.HAT_DOWN:      {joypad.ButtonDown},
	sdl.HAT_LEFTDOWN:  {joypad.ButtonDown},
	sdl.HAT_RIGHTDOWN: {joypad.ButtonDown},
	sdl.HAT_LEFT:      {joypad.ButtonLeft},
}

var playerOneKeyboardMap = map[sdl.Keycode]joypad.Button{
	sdl.K_w:         joypad.ButtonUp,
	sdl.K_s:         joypad.ButtonDown,
	sdl.K_a:         joypad.ButtonLeft,
	sdl.K_d:         joypad.ButtonRight,
	sdl.K_LSHIFT:    joypad.ButtonB,
	sdl.K_SPACE:     joypad.ButtonA,
	sdl.K_RETURN:    joypad.ButtonStart,
	sdl.K_BACKSPACE: joypad.ButtonSelect,
}

type WindowSize struct {
	Width  int
	Heigth int
}

type Window struct {
	width                 int
	heigth                int
	imagesCh              chan image.RGBA
	joypadOne             *joypad.Joypad
	joypadTwo             *joypad.Joypad
	playerOneControllerId int
	playerTwoControllerId int
}

func NewWindow(
	width int,
	heigth int,
	joypadOne *joypad.Joypad,
	joypadTwo *joypad.Joypad,
	imagesCh chan image.RGBA,
) *Window {
	return &Window{
		width:                 width,
		heigth:                heigth,
		imagesCh:              imagesCh,
		joypadOne:             joypadOne,
		joypadTwo:             joypadTwo,
		playerOneControllerId: -1,
		playerTwoControllerId: -1,
	}
}

func (w *Window) Show() {
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
	w.run(window, surface)
}

func (w *Window) run(window *sdl.Window, surface *sdl.Surface) {
eventLoop:
	for {
		select {
		case img := <-w.imagesCh:
			shouldQuit := w.handleEvents()
			if shouldQuit {
				break eventLoop
			}
			drawSurface(window, surface, img)
			window.UpdateSurface()
		case <-time.After(time.Millisecond * 10):
			shouldQuit := w.handleEvents()
			if shouldQuit {
				break eventLoop
			}
		}
	}
}

func (w *Window) handleEvents() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		shouldQuit := w.handleEvent(event)
		if shouldQuit {
			return shouldQuit
		}
	}
	return false
}

func (w *Window) handleEvent(event sdl.Event) bool {
	switch event.GetType() {
	case sdl.QUIT:
		return true
	case sdl.KEYDOWN,
		sdl.KEYUP,
		sdl.JOYBUTTONDOWN,
		sdl.JOYBUTTONUP,
		sdl.JOYHATMOTION:
		w.updateJoypadButtonsState(event)
	case sdl.JOYDEVICEADDED:
		event, ok := event.(*sdl.JoyDeviceAddedEvent)
		if !ok {
			return false
		}
		controllerId := int(event.Which)
		controller := sdl.GameControllerOpen(controllerId)
		if w.playerOneControllerId == -1 {
			w.playerOneControllerId = controllerId
		} else if w.playerTwoControllerId == -1 {
			w.playerTwoControllerId = controllerId
		}
		fmt.Printf("controller %s attached at port: %d\n", controller.Name(), controllerId)
	case sdl.JOYDEVICEREMOVED:
		event, ok := event.(*sdl.JoyDeviceRemovedEvent)
		if !ok {
			return false
		}
		controllerId := int(event.Which)
		switch controllerId {
		case w.playerOneControllerId:
			w.playerOneControllerId = -1
		case w.playerTwoControllerId:
			w.playerTwoControllerId = -1
		}
		fmt.Printf("controller at port %d disconnected\n", controllerId)
	}
	return false
}

func (w *Window) updateJoypadButtonsState(event sdl.Event) {
	switch event := event.(type) {
	case *sdl.JoyButtonEvent:
		joypd := w.getJoypad(int(event.Which))
		if joypd == nil {
			return
		}
		button, ok := joypadButtonsMap[event.Button]
		if !ok {
			return
		}
		if event.State == sdl.RELEASED {
			joypd.SetControl(button, false)
		} else {
			joypd.SetControl(button, true)
		}
	case *sdl.JoyHatEvent:
		joypd := w.getJoypad(int(event.Which))
		if joypd == nil {
			return
		}
		var buttons []joypad.Button
		var value bool
		if selectedButtons, ok := joypadDPadMap[event.Value]; ok {
			buttons = selectedButtons
			value = true
		} else {
			buttons = []joypad.Button{
				joypad.ButtonUp,
				joypad.ButtonDown,
				joypad.ButtonLeft,
				joypad.ButtonRight,
			}
			value = false
		}

		for _, button := range buttons {
			joypd.SetControl(button, value)
		}
	case *sdl.KeyboardEvent:
		pressed := event.Type == sdl.KEYDOWN
		if button, ok := playerOneKeyboardMap[event.Keysym.Sym]; ok {
			w.joypadOne.SetControl(button, pressed)
		}
	}
}

func (w *Window) getJoypad(controllerId int) *joypad.Joypad {
	switch controllerId {
	case w.playerOneControllerId:
		return w.joypadOne
	case w.playerTwoControllerId:
		return w.joypadTwo
	default:
		return nil
	}
}

func drawSurface(window *sdl.Window, surface *sdl.Surface, img image.RGBA) {
	var width, heigth = window.GetSize()
	for x := range int(width) {
		for y := range int(heigth) {
			surface.Set(x, y, img.At(x, y))
		}
	}
}
