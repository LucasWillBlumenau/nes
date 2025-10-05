package cartridge

import (
	"errors"
	"io"
	"os"
)

var ErrInvalidRomFile = errors.New("invalid rom file")

const programBanksIndex = 4
const charactersBanksIndex = 5
const firstControlByteIndex = 6
const secondControlByteIndex = 7
const ramBanksQuantityIndex = 8

const headersSize = 16
const prgBankSize = 16 * 1024

type Cartridge struct {
	UseVerticalMirroring   bool
	UseBatteryBackedRam    bool
	UseTrainer             bool
	UseFourScreenMirroring bool
	RamBanksQuantity       uint
	ProgramBanksQuantity   uint
	MapperId               uint
	Trainer                []byte
	CharacterRom           []byte
	prgRomBanks            [][]byte
}

func LoadCartridge(filePath string) (*Cartridge, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	headers := make([]byte, headersSize)

	n, err := fp.Read(headers)
	if err != nil {
		return nil, err
	}

	if n != headersSize {
		return nil, ErrInvalidRomFile
	}

	programBanksQuantity := headers[programBanksIndex]
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
		_, err = fp.Read(trainer)
		if err != nil {
			return nil, err
		}
	}

	programBanks := make([][]byte, programBanksQuantity)
	for i := range programBanks {
		programBanks[i] = make([]byte, prgBankSize)
		n, err = fp.Read(programBanks[i])
		if err != nil {
			return nil, err
		}
		if n != prgBankSize {
			return nil, ErrInvalidRomFile
		}
	}

	characterRom := make([]byte, charactersSize)
	n, err = fp.Read(characterRom)
	if err != nil {
		return nil, err
	}

	if n != int(charactersSize) {
		return nil, ErrInvalidRomFile
	}

	buf := make([]byte, 1)
	_, err = fp.Read(buf)
	if !errors.Is(err, io.EOF) {
		return nil, ErrInvalidRomFile
	}

	return &Cartridge{
		UseVerticalMirroring:   useVerticalMirroring,
		UseBatteryBackedRam:    useBatteryBackedRam,
		UseTrainer:             useTrainer,
		UseFourScreenMirroring: useFourScreenMirroring,
		RamBanksQuantity:       uint(ramBanksQuantity),
		ProgramBanksQuantity:   uint(programBanksQuantity),
		MapperId:               uint(mapperId),
		Trainer:                trainer,
		CharacterRom:           characterRom,
		prgRomBanks:            programBanks,
	}, nil
}

func (c *Cartridge) ReadPrgRom(addr uint16) uint8 {
	if c.ProgramBanksQuantity == 1 {
		if addr >= prgBankSize {
			addr -= prgBankSize
		}
		return c.prgRomBanks[0][addr]
	}

	if addr >= prgBankSize {
		return c.prgRomBanks[1][addr-prgBankSize]
	}
	return c.prgRomBanks[0][addr]
}
