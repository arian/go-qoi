package main

import (
	"flag"
	"fmt"
	"os"

	go_qoi "github.com/arian/go-qoi"
)

func main() {
	output := flag.String("output", "output.qoi", "Output qoi file to write to")
	flag.Parse()
	file := flag.Arg(0)

	fmt.Printf("file %s output %s %v\n", file, *output, flag.Args())

	if file == "" {
		flag.Usage()
		os.Exit(1)
	}

	err := go_qoi.ReadPngAndSaveImageToQoi(file, *output)
	if err != nil {
		panic(err)
	}
}
