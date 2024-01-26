package main

import (
	"flag"
	"log"
	"os"

	"github.com/vimcki/go-di-graph/internal/depcalc"
)

func main() {
	entryPoint := flag.String("entrypoint", "", "Entrypoint function")
	path := flag.String("path", "", "Path to packages")

	flag.Parse()

	result, err := depcalc.Depcalc(*entryPoint, *path)
	if err != nil {
		log.Println("Error calculating dependencies:", err)
		os.Exit(1)
	}

	os.Stdout.Write([]byte(result))
}
