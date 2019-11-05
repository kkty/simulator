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

	m := simulator.NewMachine()

	if err := m.Load(string(b)); err != nil {
		log.Fatal(err)
	}

	m.ProgramCounter, err = m.FindAddress(simulator.Label(entrypoint))

	if err != nil {
		log.Fatal(err)
	}

	stats, err := m.Run(*debug)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, "executed instructions: %v\n", stats.Executed)
	fmt.Fprintf(os.Stderr, "frequent labels: %v\n", stats.Jumps(30))
}
