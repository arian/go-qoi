package main

import (
	"flag"
	"fmt"
	"os"

	go_qoi "github.com/arian/go-qoi"
)

func main() {
	output := flag.String("output", "output.png", "Output file to write to")
	flag.Parse()
	file := flag.Arg(0)

	fmt.Printf("file %s output %s %v\n", file, *output, flag.Args())

	if file == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, err := go_qoi.Decode(f)
	if err != nil {
		panic(err)
	}

	go_qoi.SaveImageToPngFile(&img, *output)
}
