package tree

import (
	"encoding/json"
	"fmt"
	"strings"
)

type dependency struct {
	ID    int
	Name  string       `json:"name,omitempty"`
	Deps  []dependency `json:"deps,omitempty"`
	Value interface{}  `json:"value,omitempty"`
}

func Render(graph string) ([]byte, error) {
	var root dependency
	if err := json.Unmarshal([]byte(graph), &root); err != nil {
		return nil, fmt.Errorf("failed to unmarshal graph: %w", err)
	}

	fillIDs(&root, new(int))

	nodes := make(map[int]bool)

	var result []string

	generateGraphD2(root, nodes, &result)

	return []byte(strings.Join(result, "\n")), nil
}

func fillIDs(dep *dependency, id *int) {
	dep.ID = *id
	*id++
	for i := range dep.Deps {
		fillIDs(&dep.Deps[i], id)
	}
}

func generateGraphD2(dep dependency, nodes map[int]bool, result *[]string) {
	if !nodes[dep.ID] {
		name := dep.Name
		if name == "" {
			name = "aggregate"
		}
		fixedName := strings.ReplaceAll(name, "\"", "\\\"")
		if dep.Value != nil {
			rawVal := fmt.Sprintf("%v", dep.Value)
			val := strings.ReplaceAll(rawVal, "\"", "\\\"")
			*result = append(*result, fmt.Sprintf("%d: \"%v\\n%s\"\n", dep.ID, val, fixedName))
		} else {
			*result = append(*result, fmt.Sprintf("%d: \"%s\"\n", dep.ID, fixedName))
		}
		nodes[dep.ID] = true
	}
	// Recursively traverse the dependencies and print the edges
	for _, childDep := range dep.Deps {
		generateGraphD2(childDep, nodes, result)
		*result = append(*result, fmt.Sprintf("%d -> %d", dep.ID, childDep.ID))
	}
}
