package main

import (
	"bufio"
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
	output := flag.String("output", "output.qoi", "Output qoi file to write to")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to this file")

	flag.Parse()
	file := flag.Arg(0)

	fmt.Printf("file %s output %s %v\n", file, *output, flag.Args())

	if file == "" {
		flag.Usage()
		os.Exit(1)
	}

	o, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()

	img, err := go_qoi.ReadPngFile(file)
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriterSize(o, go_qoi.MaxQoiSize(img))

	encode(img, w, *cpuprofile, *memprofile)

	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func encode(img image.Image, o io.Writer, cpuprofile, memprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	err := go_qoi.Encode(o, img)
	if err != nil {
		log.Fatal(err)
	}

	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
	}
}
