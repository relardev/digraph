package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type report struct {
	Report Report `json:"report"`
}

type Report struct {
	Title  string        `json:"title"`
	Graph  GraphOrConfig `json:"graph"`
	Config GraphOrConfig `json:"config"`
}

type GraphOrConfig struct {
	Data string `json:"data"`
}

func Send(name, addr, tree, config string) error {
	report := report{
		Report: Report{
			Title: name,
			Graph: GraphOrConfig{
				Data: tree,
			},
			Config: GraphOrConfig{
				Data: config,
			},
		},
	}

	serializd, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("failed marshaling report, %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		addr,
		bytes.NewReader(serializd),
	)
	if err != nil {
		return fmt.Errorf("failed creating report request, %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed sending report, %w", err)
	}
	return nil
}
