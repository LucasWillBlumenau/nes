package debug

type DebugOption uint8

const (
	DebugOptionPauseFrameGeneration = iota
	DebugOptionExecuteNextInstruction
	DebugOptionShowInstructions
	DebugOptionDumpNametables
	DebugOptionPrintXCoordinatesOfTopScanline
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
	case DebugOptionPrintXCoordinatesOfTopScanline:
		return "DebugOptionPrintXCoordinatesOfTopScanline"
	}
	return "Invalid"

}
