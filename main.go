package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/kkty/simulator/simulator"
)

func main() {
	entrypoint := flag.String("entrypoint", "", "label name of the entrypoint")
	fileName := flag.String("file", "", "path to the assembly file")
	debug := flag.Bool("debug", false, "write debug log to stderr")
	flag.Parse()

	if *entrypoint == "" || *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	b, err := ioutil.ReadFile(*fileName)

	if err != nil {
		log.Fatal(err)
	}

	m := simulator.NewMachine(*entrypoint)

	if err := m.Load(string(b)); err != nil {
		log.Fatal(err)
	}

	m.ProgramCounter, err = m.FindAddress(simulator.Label(*entrypoint))

	if err := m.Run(*debug); err != nil {
		log.Fatal(err)
	}
}
