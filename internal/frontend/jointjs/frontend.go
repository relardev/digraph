package jointjs

import (
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"strings"
)

type Dependency struct {
	Name     string        `json:"name"`
	Deps     []*Dependency `json:"deps,omitempty"`
	Value    interface{}   `json:"value,omitempty"`
	Url      string        `json:"url,omitempty"`
	FilePath string        `json:"file_path,omitempty"`
	md5      string
}

//go:embed index.html
var htmlTemplate string

//go:embed md5.js
var md5str string

func Render(data []byte) ([]byte, error) {
	graph := &Dependency{}
	err := json.Unmarshal(data, graph)
	if err != nil {
		return nil, err
	}

	err = fillMd5(graph)
	if err != nil {
		return nil, err
	}

	fillMissingNames(graph)

	data, err = json.Marshal(graph)
	if err != nil {
		return nil, err
	}

	filled := strings.ReplaceAll(htmlTemplate, "{{GRAPH_DATA}}", string(data))
	filled = strings.ReplaceAll(filled, "{{MD5}}", md5str)

	return []byte(filled), nil
}

func fillMd5(graph *Dependency) error {
	data, err := json.Marshal(graph)
	if err != nil {
		return err
	}
	graph.md5 = md5sum(data)
	for _, dep := range graph.Deps {
		err = fillMd5(dep)
		if err != nil {
			return err
		}
	}
	return nil
}

func md5sum(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	hash := hasher.Sum(nil)

	// Convert hash to hexadecimal string
	hashString := hex.EncodeToString(hash)
	return hashString
}

func fillMissingNames(graph *Dependency) {
	hashToName := map[string]string{}
	usedNames := map[string]int{}
	fillMissingNamesRecursively(graph, hashToName, usedNames)
}

func fillMissingNamesRecursively(
	graph *Dependency,
	hashToName map[string]string,
	nameToID map[string]int,
) {
	foundName, ok := hashToName[graph.md5]
	if ok {
		graph.Name = foundName
	} else {
		if graph.Name == "" {
			graph.Name = "aggregate"
		}

		id, ok := nameToID[graph.Name]
		if ok {
			baseName := graph.Name
			graph.Name = graph.Name + "_" + strconv.Itoa(id)
			nameToID[baseName] = id + 1
		} else {
			nameToID[graph.Name] = 0
		}

		hashToName[graph.md5] = graph.Name
	}
	for _, dep := range graph.Deps {
		fillMissingNamesRecursively(dep, hashToName, nameToID)
	}
}
