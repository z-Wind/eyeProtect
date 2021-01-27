package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

const (
	screenWidth  = 800
	screenHeight = 600
	dpi          = 90
	size         = 56
)

var counter int

func init() {
	flag.IntVar(&counter, "t", 20, "總秒數")
}

func main() {
	flag.Parse()

	ch := make(chan int)

	go func() {
		for i := range ch {
			fmt.Println(i)
		}
		fmt.Printf("close\n")
		os.Exit(0)
	}()

	switch os := runtime.GOOS; os {
	case "windows":
		ebitenMain(counter, ch)
	case "linux":
		oakMain(counter, ch)
	default:
		log.Fatalf("Unsupported OS: %v", os)
	}

}
