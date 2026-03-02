package cartridge

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type MirroringType uint8

const (
	VerticalMirroring MirroringType = iota
	HorizontalMirroring
	FourScreenMirroring
)

var ErrInvalidRomFile = errors.New("invalid rom file")
var ErrUnimplementedMapper = errors.New("unimplemented mapper")

const (
	programBanksIndex      = 4
	charactersBanksIndex   = 5
	firstControlByteIndex  = 6
	secondControlByteIndex = 7
	ramBanksQuantityIndex  = 8

	headersSize = 16
	prgBankSize = 16 * 1024
	chrBankSize = 8 * 1024
)

type characterMemory interface {
	Read(addr uint16) uint8
	Write(addr uint16, data uint8)
}

type characterRam []byte

func (r characterRam) Read(addr uint16) uint8 {
	return (r)[addr]
}

func (r characterRam) Write(addr uint16, data uint8) {
	(r)[addr] = data
}

type characterRom []byte

func (r characterRom) Read(addr uint16) uint8 {
	return (r)[addr]
}

func (r characterRom) Write(addr uint16, data uint8) {
}

type cartridgeHeaders struct {
	ProgramBanksQuantity   int
	ProgramRomSize         int
	CharacterBanksQuantity int
	CharacterRomSize       int
	Mirroring              MirroringType
	UseBatteryBackedRam    bool
	UseCharacterRam        bool
	UseTrainer             bool
	RamBanksQuantity       int
	MapperId               int
}

type cartridgeRom struct {
	Character characterMemory
	Program   []byte
	Trainers  []byte
}

type Cartridge struct {
	headers cartridgeHeaders
	mapper  mapper
}

func LoadCartridgeFromRom(filePath string) (*Cartridge, error) {
	fp, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	headers, err := readHeaders(fp)
	if err != nil {
		return nil, err
	}

	rom, err := readRom(fp, headers)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1)
	_, err = fp.Read(buf)
	if !errors.Is(err, io.EOF) {
		return nil, ErrInvalidRomFile
	}

	var createMapper createMapperFn = nil
	if len(mappers) > headers.MapperId {
		createMapper = mappers[headers.MapperId]
	}
	if createMapper == nil {
		return nil, fmt.Errorf("%w: mapper %d not implemented", ErrUnimplementedMapper, headers.MapperId)
	}
	mapper := createMapper(rom, headers)
	return &Cartridge{
		headers: *headers,
		mapper:  mapper,
	}, nil
}

func readHeaders(reader io.Reader) (*cartridgeHeaders, error) {
	headers := make([]byte, headersSize)
	n, err := reader.Read(headers)
	if err != nil {
		return nil, err
	}

	if n != headersSize {
		return nil, ErrInvalidRomFile
	}

	programBanksQuantity := int(headers[programBanksIndex])
	charactersBanks := int(headers[charactersBanksIndex])
	useCharacterRom := false
	if charactersBanks == 0 {
		charactersBanks = 1
		useCharacterRom = true
	}

	charactersSize := charactersBanks * chrBankSize
	programSize := programBanksQuantity * prgBankSize
	firstControlByte := headers[firstControlByteIndex]
	secondControlByte := headers[secondControlByteIndex]
	ramBanksQuantity := int(headers[ramBanksQuantityIndex])

	useVerticalMirroring := (firstControlByte & 0b1) == 1
	useFourScreenMirroring := (firstControlByte & 0b1000) == 1
	mirroring := HorizontalMirroring
	if useFourScreenMirroring {
		mirroring = FourScreenMirroring
	} else if useVerticalMirroring {
		mirroring = VerticalMirroring
	}

	useBatteryBackedRam := (firstControlByte & 0b10) == 1
	mapperId := int((secondControlByte & 0b11110000) | (firstControlByte >> 4))

	useTrainer := (firstControlByte & 0b100) == 1

	return &cartridgeHeaders{
		ProgramBanksQuantity:   programBanksQuantity,
		ProgramRomSize:         programSize,
		CharacterBanksQuantity: charactersBanks,
		CharacterRomSize:       charactersSize,
		Mirroring:              mirroring,
		UseBatteryBackedRam:    useBatteryBackedRam,
		UseCharacterRam:        useCharacterRom,
		UseTrainer:             useTrainer,
		RamBanksQuantity:       ramBanksQuantity,
		MapperId:               mapperId,
	}, nil

}

func readRom(reader io.Reader, headers *cartridgeHeaders) (*cartridgeRom, error) {
	var trainer []byte
	if headers.UseTrainer {
		trainer = make([]byte, 512)
		_, err := reader.Read(trainer)
		if err != nil {
			return nil, err
		}
	}

	prgRom := make([]byte, headers.ProgramRomSize)
	n, err := reader.Read(prgRom)
	if err != nil {
		return nil, err
	}

	if n != headers.ProgramRomSize {
		return nil, ErrInvalidRomFile
	}

	var chrRomData = make([]byte, headers.CharacterRomSize)
	var chrRom characterMemory
	if headers.UseCharacterRam {
		chrRom = characterRam(chrRomData)
	} else {
		n, err = reader.Read(chrRomData)
		if err != nil {
			return nil, err
		}
		if n != headers.CharacterRomSize {
			return nil, ErrInvalidRomFile
		}
		chrRom = characterRom(chrRomData)
	}

	return &cartridgeRom{
		Character: chrRom,
		Program:   prgRom,
		Trainers:  trainer,
	}, nil
}

func (c *Cartridge) Mirroring() MirroringType {
	return c.mapper.Mirroring()
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
