package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/png"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/amitu/gutils"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	// quality := flag.Float64("quality", 0.7, "Quality of compression.")
	benchmark := flag.Bool("benchmark", false, "Benchmark it.")
	iterations := flag.Int("iterations", 100, "Iterations for benchmark")
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	filename := flag.Arg(0)

	if *benchmark {
		start := time.Now()
		b := gutils.Screenshot().Bounds()
		buf := make([]byte, b.Dx()*b.Dy()*4)
		w := bytes.NewBuffer(buf)
		for i := 0; i < *iterations; i++ {
			m := gutils.Screenshot()
			if filename != "" {
				// f, err := os.Create(filename)
				// if err != nil {
				// 	fmt.Println("Benchmark failed:", err)
				// 	return
				// }
				start := time.Now()
				// jpeg.Encode(w, m, nil)
				png.Encode(w, m)
				fmt.Println(time.Since(start))
				// f.Close()
			}
		}
		delta := time.Now().Sub(start)
		usedDisk := filename != ""
		fmt.Printf(
			"Result: iterations=%d, time=%s, fps=%d, disk=%t.\n", *iterations,
			delta, *iterations*1e9/int(delta), usedDisk,
		)
	} else {
		if filename == "" {
			fmt.Println("Filename required.")
			return
		}

		m := gutils.Screenshot()
		f, err := os.Create(filename)
		if err != nil {
			fmt.Println("Benchmark failed:", err)
			return
		}
		png.Encode(f, m)
		f.Close()
	}
}
