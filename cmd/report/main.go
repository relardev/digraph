package main

import (
	"os"
	"time"

	"github.com/vimcki/go-di-graph/internal/report"
)

func main() {
	name := "freya"
	graphData, err := os.ReadFile("_projects/" + name + "/enhanced.json")
	if err != nil {
		panic(err)
	}

	configData, err := os.ReadFile("_projects/" + name + "/config.json")
	if err != nil {
		panic(err)
	}

	// run 100 times
	for i := 0; i < 100; i++ {
		go report.Send(
			name,
			"http://localhost:4000/api/reports",
			string(graphData),
			string(configData),
		)
	}
	time.Sleep(10 * time.Second)
}
