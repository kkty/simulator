package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kkty/simulator/pkg/simulator"
)

func main() {
	mappedFile := flag.String("mapped-file", "", "specifies the file to map on memory")
	mappedAddress := flag.Int("mapped-address", 0, "specifies where to map the input on memory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: simulator [OPTIONS] FILENAME ENTRYPOINT\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	fileName := flag.Arg(0)
	entrypoint := flag.Arg(1)

	if entrypoint == "" || fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	var mappedData []byte
	if len(*mappedFile) > 0 {
		var err error
		mappedData, err = ioutil.ReadFile(*mappedFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	m := simulator.NewMachine()

	if err := m.Load(string(b), mappedData, int32(*mappedAddress)); err != nil {
		log.Fatal(err)
	}

	m.ProgramCounter, err = m.FindAddress(simulator.Label(entrypoint))

	if err != nil {
		log.Fatal(err)
	}

	executed, err := m.Run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "executed instructions: %v\n", executed)
}
