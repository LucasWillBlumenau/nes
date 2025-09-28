package cartridge

import (
	"errors"
	"io"
)

var ErrInvalidRomFile = errors.New("invalid rom file")

const programBanksIndex = 4
const charactersBanksIndex = 5
const firstControlByteIndex = 6
const secondControlByteIndex = 7
const ramBanksQuantityIndex = 8

const headersSize = 16

type Cartridge struct {
	UseVerticalMirroring   bool
	UseBatteryBackedRam    bool
	UseTrainer             bool
	UseFourScreenMirroring bool
	RamBanksQuantity       uint
	ProgramBanks           uint
	MapperId               uint
	Trainer                []byte
	ProgramRom             []byte
	CharacterRom           []byte
}

func LoadCartridgeFromReader(r io.ReadCloser) (*Cartridge, error) {
	defer r.Close()

	headers := make([]byte, headersSize)

	n, err := r.Read(headers)
	if err != nil {
		return nil, err
	}

	if n != headersSize {
		return nil, ErrInvalidRomFile
	}

	programBanks := headers[programBanksIndex]
	programSize := uint(programBanks) * 16 * 1024
	charactersBanks := headers[charactersBanksIndex]
	charactersSize := uint(charactersBanks) * 8 * 1024
	firstControlByte := headers[firstControlByteIndex]
	secondControlByte := headers[secondControlByteIndex]
	ramBanksQuantity := headers[ramBanksQuantityIndex]

	useVerticalMirroring := (firstControlByte & 0b1) == 1
	useBatteryBackedRam := (firstControlByte & 0b10) == 1
	useFourScreenMirroring := (firstControlByte & 0b1000) == 1

	mapperId := (secondControlByte & 0b11110000) | (firstControlByte >> 4)

	useTrainer := (firstControlByte & 0b100) == 1
	var trainer []byte
	if useTrainer {
		trainer = make([]byte, 512)
		_, err = r.Read(trainer)
		if err != nil {
			return nil, err
		}
	}

	programRom := make([]byte, programSize)
	n, err = r.Read(programRom)
	if err != nil {
		return nil, err
	}

	if n != int(programSize) {
		return nil, ErrInvalidRomFile
	}

	characterRom := make([]byte, charactersSize)
	n, err = r.Read(characterRom)
	if err != nil {
		return nil, err
	}

	if n != int(charactersSize) {
		return nil, ErrInvalidRomFile
	}

	buf := make([]byte, 1)
	_, err = r.Read(buf)
	if !errors.Is(err, io.EOF) {
		return nil, ErrInvalidRomFile
	}

	return &Cartridge{
		UseVerticalMirroring:   useVerticalMirroring,
		UseBatteryBackedRam:    useBatteryBackedRam,
		UseTrainer:             useTrainer,
		UseFourScreenMirroring: useFourScreenMirroring,
		RamBanksQuantity:       uint(ramBanksQuantity),
		ProgramBanks:           uint(programBanks),
		MapperId:               uint(mapperId),
		Trainer:                trainer,
		ProgramRom:             programRom,
		CharacterRom:           characterRom,
	}, nil

}
