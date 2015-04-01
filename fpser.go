package gutils

import "time"

// FPSer is used to keep things ticking at predefined FPS.
type FPSer struct {
	fps    float32
	last   time.Time
	target time.Duration
}

// Tick sleeps if required to keep the FPS.
func (f *FPSer) Tick() {
	delta := time.Since(f.last)
	if delta < f.target {
		time.Sleep(f.target - delta)
	}
	f.last = time.Now()
}

// NewFPSer can be used to construct FPSser
func NewFPSer(fps float32) *FPSer {
	target := time.Duration(float32(time.Second) / fps)
	return &FPSer{fps, time.Now().Add(-target), target}
}
