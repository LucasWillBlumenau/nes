package interrupt

type Interrupt uint8

const (
	Reset                = 0
	NonMaskableInterrupt = 1
	Irq                  = 2
)

func (i Interrupt) String() string {
	switch i {
	case Reset:
		return "Reset"
	case NonMaskableInterrupt:
		return "NonMaskableInterrupt"
	case Irq:
		return "Irq"
	}
	return ""
}

type interruptSignal struct {
	wasTriggered bool
	value        Interrupt
}

func (s *interruptSignal) Read() (Interrupt, bool) {
	wasTriggered := s.wasTriggered
	s.wasTriggered = false
	return s.value, wasTriggered
}

func (s *interruptSignal) Send(interrupt Interrupt) {
	s.value = interrupt
	s.wasTriggered = true

}

var InterruptSignal = interruptSignal{wasTriggered: false}
