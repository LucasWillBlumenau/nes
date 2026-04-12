package debug

type DebugOption uint8

const (
	DebugOptionPauseFrameGeneration = iota
	DebugOptionExecuteNextInstruction
	DebugOptionShowInstructions
	DebugOptionDumpNametables
)

func (d DebugOption) String() string {
	switch d {
	case DebugOptionPauseFrameGeneration:
		return "DebugOptionPauseFrameGeneration"
	case DebugOptionExecuteNextInstruction:
		return "DebugOptionExecuteNextInstruction"
	case DebugOptionShowInstructions:
		return "DebugOptionShowInstructions"
	case DebugOptionDumpNametables:
		return "DebugOptionDumpNametables"
	}
	return "Invalid"

}
