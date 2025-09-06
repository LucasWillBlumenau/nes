package main

import (
	"log"
	"os"

	"github.com/LucasWillBlumenau/nes/cartridge"
)

func main() {

	fp, err := os.Open("rom.nes")
	if err != nil {
		log.Fatalln(err)
	}

	_, err = cartridge.LoadCartridgeFromReader(fp)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Got to the end without errors")

}
