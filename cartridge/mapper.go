package cartridge

type mapper interface {
	ReadPrg(addr uint16) uint8
	WritePrg(addr uint16, data uint8)
	ReadChr(addr uint16) uint8
	WriteChr(addr uint16, data uint8)
}
