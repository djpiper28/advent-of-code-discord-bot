package main

import (
	"gonum.org/v1/plot/vg"
	"image/color"
	"log"
	"sync"
	"time"
)

// / Plots are stored on the disk in a scratch file whilst they are being made
const PLOT_SCRATCH_FILE = "tmp.plot_scratch_file.png"
const PLOT_SIZE = 8 * vg.Inch
const BG_COLOUR = 0x36393F
const TEXT_COLOUR = 0xFFFFFF

func TimeToPlot(t time.Time) float64 {
	return float64(t.Day()) + float64(t.Hour()/24) + float64(t.Minute())/24/60
}

func HexToRGB(hex uint) color.Color {
	return color.RGBA{R: uint8(hex & 0xFF0000 >> 16),
		G: uint8(hex & 0xFF00 >> 8),
		B: uint8(hex & 0xFF),
		A: 0xFF}
}

var plotScratchFileLock sync.Mutex

func LockAndPlot(f func()) {
	plotScratchFileLock.Lock()
	defer func() {
		plotScratchFileLock.Unlock()

		err := recover()
		if err != nil {
			log.Print(err)
		}
	}()

	f()
}
