package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/kkty/simulator/simulator"
)

func main() {
	b, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal(err)
	}

	m := simulator.NewMachine(1000)

	if err := m.Load(string(b)); err != nil {
		log.Fatal(err)
	}

	m.ProgramCounter, err = m.FindAddress("min_caml_start")

	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}
