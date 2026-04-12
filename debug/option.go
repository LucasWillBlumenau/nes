package debug

type DebugOption uint8

const (
	DebugOptionPauseFrameGeneration = iota
	DebugOptionExecuteNextInstruction
	DebugOptionShowInstructions
)

func (c DebugOption) String() string {
	switch c {
	case DebugOptionPauseFrameGeneration:
		return "DebugOptionPauseFrameGeneration"
	case DebugOptionExecuteNextInstruction:
		return "DebugOptionExecuteNextInstruction"
	}
	return "Invalid"

}
