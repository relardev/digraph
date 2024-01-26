package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/vimcki/go-di-graph/internal/frontend/d2"
)

func main() {
	graphFilePath := flag.String("graph_path", "", "Path to the graph file")

	flag.Parse()

	data, err := os.ReadFile(*graphFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	result, err := d2.Render(string(data))
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	os.Stdout.Write([]byte(result))
}
