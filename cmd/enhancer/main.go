package main

import (
	"flag"
	"log"
	"os"

	"github.com/vimcki/go-di-graph/internal/enhancer"
)

func main() {
	configPath := flag.String("config_path", "", "Path to config file")
	treePath := flag.String("tree_path", "", "Path to tree file")
	projectName := flag.String("project_name", "", "Name of the project")
	baseUrl := flag.String("base_url", "", "Base url of the project")

	flag.Parse()

	data, err := os.ReadFile(*treePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	e, err := enhancer.New(
		*configPath,
		enhancer.WithMetadata(*projectName, *baseUrl),
	)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	rusult, err := e.Enhance(string(data))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Stdout.Write([]byte(rusult))
}
