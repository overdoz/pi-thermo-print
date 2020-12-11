package main

import (
	"github.com/disintegration/imaging"
	"log"
	"os/exec"
	"strings"
)

// edit immage
func AdjustImage(file string) error {
	src, _ := imaging.Open(file)
	src = imaging.Resize(src, 800, 0, imaging.Lanczos)
	src = imaging.AdjustBrightness(src, 30)
	src = imaging.Grayscale(src)
	src = imaging.AdjustContrast(src, -20)
	src = imaging.AdjustGamma(src, 0.75)
	return imaging.Save(src, FilenamePng)
}

// file should have the ending .png
func PrintImage(file string) {
	sh := "lp " + file + " -d LKT"

	args := strings.Split(sh, " ")

	cmd := exec.Command(args[0], args[1:]...)

	_, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(err)
	}
}
