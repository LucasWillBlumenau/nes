package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	filePath := readCliArgs()
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	lastIndex := len(content) - 1
	for _, b := range content[:lastIndex] {
		buf.Write(fmt.Appendf(nil, "%02x ", b))
	}
	buf.Write(fmt.Appendf(nil, "%02x", content[lastIndex]))

	fp, err := os.Create("dump.nes")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	fp.Write(buf.Bytes())
}

func readCliArgs() string {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("the program only supports a rom path as argument")
	}
	return args[0]
}
