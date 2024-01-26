package d2

import "github.com/vimcki/go-di-graph/internal/frontend/d2/dag"

func Render(graph string) ([]byte, error) {
	return dag.Render(graph)
}
