package main

import (
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"runtime/pprof"

	go_qoi "github.com/arian/go-qoi"
)

func main() {
	output := flag.String("output", "output.png", "Output file to write to")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	file := flag.Arg(0)

	fmt.Printf("file %s output %s %v\n", file, *output, flag.Args())

	if file == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img := decode(f, *cpuprofile)

	go_qoi.SaveImageToPngFile(&img, *output)
}

func decode(f io.Reader, cpuprofile string) image.Image {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	img, err := go_qoi.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	return img
}
