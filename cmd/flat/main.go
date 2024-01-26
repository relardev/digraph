package main

import (
	"flag"
	"log"
	"os"

	"github.com/vimcki/go-di-graph/internal/flatten"
)

func main() {
	config := flag.String("config", "", "Path to config file")
	entryPoint := flag.String("entrypoint", "", "Entrypoint function")
	basePath := flag.String("basepath", "", "Base path of the project")
	buildPackage := flag.String("buildpackage", "", "Package to build")
	flatPackage := flag.String("flatpackage", "", "Package to flatten")

	flag.Parse()

	err := flatten.Flatten(*basePath, *buildPackage, *flatPackage, *entryPoint, *config)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
