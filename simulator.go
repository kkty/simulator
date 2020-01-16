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
	native := flag.Bool("native", false, "use native float operation")

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

	m := simulator.NewMachine()

	if err := m.Load(string(b)); err != nil {
		log.Fatal(err)
	}

	m.ProgramCounter, err = m.FindAddress(simulator.Label(entrypoint))

	if err != nil {
		log.Fatal(err)
	}

	executed, err := m.Run(*native)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "executed instructions: %v\n", executed)
}
