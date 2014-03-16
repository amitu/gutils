package main

import (
	"flag"
	"io/ioutil"
	"fmt"
	"time"
	"github.com/amitu/gutils/macshot"
)

func main() {
	quality := flag.Float64("quality", 0.7, "Quality of compression.")
	benchmark := flag.Bool("benchmark", false, "Benchmark it.")
	iterations := flag.Int("iterations", 100, "Iterations for benchmark")

	flag.Parse()
	filename := flag.Arg(0)

	if *benchmark {
		start := time.Now()
		for i := 0; i < *iterations; i++ {
			jpeg, err := macshot.ScreenShot(*quality)
			if err != nil {
				fmt.Println("Benchmark failed:", err)
				return
			}
			if filename != "" {
				err = ioutil.WriteFile(filename, jpeg, 0644)
				if err != nil {
					fmt.Println("Benchmark failed:", err)
					return
				}
			}
		}
		delta := time.Now().Sub(start)
		usedDisk := filename != ""
		fmt.Printf(
			"Result: iterations=%d, time=%s, fps=%d, disk=%t.\n", *iterations,
			delta, *iterations * 1e9 / int(delta), usedDisk,
		)
	} else {
		jpeg, err := macshot.ScreenShot(*quality)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = ioutil.WriteFile(filename, jpeg, 0644)
		if err != nil {
			fmt.Println(err)
		}
	}
}