package main

import (
	"bufio"
	"flag"
	"log"
	"os"

	"github.com/pavlovskyive/architecture-lab-4/commands"
	"github.com/pavlovskyive/architecture-lab-4/engine"
)

var (
	inputFile = flag.String("f", "", "Input file")
)

func main() {
	flag.Parse()

	eventLoop := new(engine.EventLoop)
	eventLoop.Start()

	if *inputFile == "" {
		log.Fatalln("Input file not specified. Specify input file after -f flag.")
	} else {
		println("Read from file: " + *inputFile)
	}

	if input, err := os.Open(*inputFile); err == nil {
		defer input.Close()

		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			commandLine := scanner.Text()
			cmd := commands.Parse(commandLine)
			eventLoop.Post(cmd)
		}
	} else if err != nil {
		log.Fatalln(err.Error())
	}

	eventLoop.AwaitFinish()
}
