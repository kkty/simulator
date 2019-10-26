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
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: simulator [OPTIONS] FILENAME ENTRYPOINT\n")
		flag.PrintDefaults()
	}

	debug := flag.Bool("debug", false, "write debug log to stderr")

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

	m := simulator.NewMachine(entrypoint)

	if err := m.Load(string(b)); err != nil {
		log.Fatal(err)
	}

	m.ProgramCounter, err = m.FindAddress(simulator.Label(entrypoint))

	if err != nil {
		log.Fatal(err)
	}

	if err := m.Run(*debug); err != nil {
		log.Fatal(err)
	}
}
