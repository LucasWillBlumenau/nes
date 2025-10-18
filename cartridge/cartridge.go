package cartridge

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var ErrInvalidRomFile = errors.New("invalid rom file")
var ErrUnimplementedMapper = errors.New("unimplemented mapper")

const programBanksIndex = 4
const charactersBanksIndex = 5
const firstControlByteIndex = 6
const secondControlByteIndex = 7
const ramBanksQuantityIndex = 8

const headersSize = 16
const prgBankSize = 16 * 1024

type Cartridge struct {
	useVerticalMirroring   bool
	useBatteryBackedRam    bool
	useTrainer             bool
	useFourScreenMirroring bool
	ramBanksQuantity       int
	mapperId               int
	Trainer                []byte
	mapper                 mapper
}

func LoadCartridgeFromRom(filePath string) (*Cartridge, error) {
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

	prgRomSize := int(programBanksQuantity) * prgBankSize
	programBanks := make([]byte, prgRomSize)
	n, err = fp.Read(programBanks)
	if err != nil {
		return nil, err
	}

	if n != prgRomSize {
		return nil, ErrInvalidRomFile
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

	if mapperId != 0 {
		return nil, fmt.Errorf("%w: mapper %d not implemented", ErrUnimplementedMapper, mapperId)
	}

	mapper := &nrom{
		BanksQuantity: int(programBanksQuantity),
		PrgRom:        programBanks,
		ChrRom:        characterRom,
	}

	return &Cartridge{
		useVerticalMirroring:   useVerticalMirroring,
		useBatteryBackedRam:    useBatteryBackedRam,
		useTrainer:             useTrainer,
		useFourScreenMirroring: useFourScreenMirroring,
		ramBanksQuantity:       int(ramBanksQuantity),
		mapperId:               int(mapperId),
		Trainer:                trainer,
		mapper:                 mapper,
	}, nil
}

func (c *Cartridge) UseVerticalMirroring() bool {
	return c.useVerticalMirroring
}

func (c *Cartridge) ReadPrgRom(addr uint16) uint8 {
	return c.mapper.ReadPrg(addr)
}

func (c *Cartridge) WritePrgRom(addr uint16, data uint8) {
	c.mapper.WritePrg(addr, data)
}

func (c *Cartridge) ReadChrRom(addr uint16) uint8 {
	return c.mapper.ReadChr(addr)
}

func (c *Cartridge) WriteChrRom(addr uint16, data uint8) {
	c.mapper.WriteChr(addr, data)
}
