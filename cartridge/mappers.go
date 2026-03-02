package cartridge

const lastMapperId = 232

type createMapperFn func(rom *cartridgeRom, headers *cartridgeHeaders) mapper

var mappers = [lastMapperId + 1]createMapperFn{
	0: newINES0,
	2: newINES2,
}
