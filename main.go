package main

import (
	"log"
	"runtime"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
	dpi          = 90
	size         = 48
)

func main() {
	switch os := runtime.GOOS; os {
	case "windows":
		ebitenMain()
	case "linux":
		oakMain()
	default:
		log.Fatalf("Unsupported OS: %v", os)
	}
}
