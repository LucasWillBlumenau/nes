package interruption

type Interruption uint8

const (
	Reset                = 0
	NonMaskableInterrupt = 1
	Irq                  = 2
)

var InterruptionHandler = make(chan Interruption)
