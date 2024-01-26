package enhancer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/maja42/goval"
)

type Dependency struct {
	Name         string        `json:"name"`
	Deps         []*Dependency `json:"deps"`
	ImportedFrom string        `json:"imported_from"`
	Value        interface{}
	Url          string
	FilePath     string
}

type OutputDependency struct {
	Name     string              `json:"name"`
	Deps     []*OutputDependency `json:"deps,omitempty"`
	Value    interface{}         `json:"value,omitempty"`
	Url      string              `json:"url,omitempty"`
	FilePath string              `json:"file_path,omitempty"`
}

type Enhancer struct {
	config      map[string]interface{}
	baseUrl     string
	projectName string
}

type Option func(*Enhancer)

func New(configPath string, options ...Option) (*Enhancer, error) {
	// Load config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config map[string]interface{}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	e := &Enhancer{
		config: config,
	}

	for _, option := range options {
		option(e)
	}

	return e, nil
}

func WithMetadata(projectName, baseUrl string) Option {
	return func(e *Enhancer) {
		e.baseUrl = baseUrl
		e.projectName = projectName
	}
}

func (e *Enhancer) Enhance(treeData string) (string, error) {
	var tree Dependency

	err := json.Unmarshal([]byte(treeData), &tree)
	if err != nil {
		return "", fmt.Errorf("failed to parse tree file: %w", err)
	}

	err = e.enhanceTree(&tree)
	if err != nil {
		return "", fmt.Errorf("failed to enhance tree: %w", err)
	}

	outputTree := toOutputTree(&tree)

	bytes, err := json.Marshal(outputTree)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tree: %w", err)
	}

	return string(bytes), nil
}

func (e *Enhancer) enhanceTree(dep *Dependency) error {
	if strings.HasPrefix(dep.Name, "cfg.") {
		ctx := map[string]interface{}{
			"cfg": e.config,
		}
		evaluator := goval.NewEvaluator()

		var err error

		dep.Value, err = evaluator.Evaluate(dep.Name, ctx, nil)
		if err != nil {
			dep.Value = dep.Name + " (error: " + err.Error() + ")"
		}
	}

	if dep.ImportedFrom != "" && e.projectName != "" {
		if strings.Contains(dep.ImportedFrom, "/"+e.projectName+"/") {
			split := strings.Split(dep.ImportedFrom, "/"+e.projectName+"/")
			dep.FilePath = split[len(split)-1]
			if e.baseUrl != "" {
				dep.Url = e.baseUrl + dep.FilePath
			}
		} else {
			if strings.Contains(dep.ImportedFrom, ".") {
				dep.Url = "https://" + dep.ImportedFrom
			}
		}
	}

	for _, depChild := range dep.Deps {
		err := e.enhanceTree(depChild)
		if err != nil {
			return err
		}
	}

	return nil
}

func toOutputTree(dep *Dependency) *OutputDependency {
	out := &OutputDependency{
		Name:     dep.Name,
		Value:    dep.Value,
		Url:      dep.Url,
		FilePath: dep.FilePath,
	}

	for _, depChild := range dep.Deps {
		out.Deps = append(out.Deps, toOutputTree(depChild))
	}

	return out
}
