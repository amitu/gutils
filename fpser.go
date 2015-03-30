package gutils

import (
	"fmt"
	"time"
)

// FPSer is used to keep things ticking at predefined FPS.
type FPSer struct {
	fps    float32
	last   time.Time
	target time.Duration
}

// Tick sleeps if required to keep the FPS.
func (f *FPSer) Tick() {
	delta := time.Since(f.last)
	// 1m1.221363664s 200ms -1m1.021363664s 1e+09 5 ft
	fmt.Println(delta, f.target, f.target-delta, float32(time.Second), f.fps, "ft")
	if delta < f.target {
		time.Sleep(f.target - delta)
	}
	f.last = time.Now()
}

// NewFPSer can be used to construct FPSser
func NewFPSer(fps float32) *FPSer {
	target := time.Duration(float32(time.Second) / fps)
	fmt.Println(target)
	return &FPSer{fps, time.Now().Add(-time.Minute), target}
}
