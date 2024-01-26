package digraph

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sync"

	"github.com/vimcki/go-di-graph/internal/depcalc"
	"github.com/vimcki/go-di-graph/internal/encoder"
	"github.com/vimcki/go-di-graph/internal/enhancer"
	"github.com/vimcki/go-di-graph/internal/flatten"
	"github.com/vimcki/go-di-graph/internal/frontend/d2"
	"github.com/vimcki/go-di-graph/internal/frontend/jointjs"
	"github.com/vimcki/go-di-graph/internal/report"
)

type reportConfig struct {
	addr string
	id   string
}

type Digraph struct {
	marshaler    func(interface{}) ([]byte, error)
	config       string
	dir          string
	graph        []byte
	contentType  string
	entrypoint   string
	buildFS      fs.FS
	err          error
	lock         sync.RWMutex
	render       func([]byte) ([]byte, string, error)
	blocking     bool
	repoName     string
	baseUrl      string
	reportConfig reportConfig
}

func New(
	config interface{},
	entrypoint string,
	buildFS fs.FS,
	options ...func(*Digraph),
) (*Digraph, error) {
	graph := &Digraph{
		entrypoint: entrypoint,
		buildFS:    buildFS,
		marshaler:  json.Marshal,
	}

	graph.render = graph.htmlRender

	for _, option := range options {
		option(graph)
	}

	cfgString, err := graph.marshaler(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	graph.config = string(cfgString)

	return graph, nil
}

func WithSendReport(id, addr string) func(*Digraph) {
	return func(d *Digraph) {
		d.reportConfig = reportConfig{
			id:   id,
			addr: addr,
		}
	}
}

func WithCustomMarshal() func(*Digraph) {
	return func(d *Digraph) {
		d.marshaler = encoder.Marshal
	}
}

func WithD2Render() func(*Digraph) {
	return func(d *Digraph) {
		d.render = d.renderD2
	}
}

func WithNoRender() func(*Digraph) {
	return func(d *Digraph) {
		d.render = func(tree []byte) ([]byte, string, error) {
			return tree, "application/json", nil
		}
	}
}

func WithBlockingHandler() func(*Digraph) {
	return func(d *Digraph) {
		d.blocking = true
	}
}

func WithMetadata(repositoryName, baseUrl string) func(*Digraph) {
	return func(d *Digraph) {
		d.repoName = repositoryName
		d.baseUrl = baseUrl
	}
}

func (d *Digraph) Close() error {
	return os.RemoveAll(d.dir)
}

func (d *Digraph) Handler() func() ([]byte, string) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := d.buildGraph()
		if err != nil {
			d.lock.Lock()

			defer d.lock.Unlock()

			d.err = err
		}
	}()

	if d.blocking {
		wg.Wait()
	}

	return d.handle
}

func (d *Digraph) buildGraph() error {
	dir, err := os.MkdirTemp("/tmp/", "digraph")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}

	d.dir = dir

	err = d.unpackBuildFS()
	if err != nil {
		return fmt.Errorf("failed to unpack build fs: %w", err)
	}

	graph, err := d.buildTree()
	if err != nil {
		return fmt.Errorf("failed to build graph: %w", err)
	}

	rendered, contentType, err := d.render(graph)
	if err != nil {
		return fmt.Errorf("failed to render graph: %w", err)
	}

	if d.reportConfig.addr != "" {
		err = report.Send(d.reportConfig.id, d.reportConfig.addr, string(graph), d.config)
		if err != nil {
			return fmt.Errorf("failed to send report, %w", err)
		}
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	d.graph = rendered
	d.contentType = contentType

	return nil
}

func (d *Digraph) handle() ([]byte, string) {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if d.err != nil {
		return []byte(d.err.Error()), "text/plain"
	}

	graph := d.graph

	if graph == nil {
		return []byte("building..."), "text/plain"
	}

	return graph, d.contentType
}

func (d *Digraph) buildTree() ([]byte, error) {
	err := os.Mkdir(d.dir+"/flat", 0o755)
	if err != nil {
		return nil, fmt.Errorf("failed to create build dir: %w", err)
	}

	configPath := path.Join(d.dir, "/config.json")

	err = os.WriteFile(configPath, []byte(d.config), 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to write config file: %w", err)
	}

	err = flatten.Flatten(d.dir, "build", "flat", d.entrypoint, configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to flatten: %w", err)
	}

	deptree, err := depcalc.Depcalc(d.entrypoint, path.Join(d.dir, "/flat"))
	if err != nil {
		return nil, fmt.Errorf("failed to calculate dependencies: %w", err)
	}

	err = os.WriteFile(path.Join(d.dir, "/deptree.json"), []byte(deptree), 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to write deptree file: %w", err)
	}

	e, err := enhancer.New(
		configPath,
		enhancer.WithMetadata(d.repoName, d.baseUrl),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create enhancer: %w", err)
	}

	enhanced, err := e.Enhance(deptree)
	if err != nil {
		return nil, fmt.Errorf("failed to enhance tree: %w", err)
	}

	err = os.WriteFile(path.Join(d.dir, "/enhanced.json"), []byte(enhanced), 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to write enhanced file: %w", err)
	}

	return []byte(enhanced), nil
}

func (d *Digraph) renderD2(tree []byte) ([]byte, string, error) {
	_, err := exec.LookPath("d2")
	if err != nil {
		return nil, "", fmt.Errorf("d2 binary is not available: %w", err)
	}

	d2Graph, err := d2.Render(string(tree))
	if err != nil {
		return nil, "", fmt.Errorf("failed to render d2 graph: %w", err)
	}

	renderPath := path.Join(d.dir, "/render.d2")

	err = os.WriteFile(renderPath, d2Graph, 0o644)
	if err != nil {
		return nil, "", fmt.Errorf("failed to write d2 file: %w", err)
	}

	svgPath := path.Join(d.dir, "/render.svg")

	cmd := exec.Command("d2", "--layout=elk", renderPath, svgPath)

	err = cmd.Run()
	if err != nil {
		return nil, "", fmt.Errorf("failed to run d2: %w", err)
	}

	bytes, err := os.ReadFile(svgPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file: %w", err)
	}

	return bytes, "image/svg+xml", nil
}

func (d *Digraph) htmlRender(tree []byte) ([]byte, string, error) {
	bytes, err := jointjs.Render(tree)
	if err != nil {
		return nil, "", fmt.Errorf("failed to render html: %w", err)
	}
	return bytes, "text/html", nil
}

func (d *Digraph) unpackBuildFS() error {
	err := os.Mkdir(d.dir+"/build", 0o755)
	if err != nil {
		return fmt.Errorf("failed to create build dir: %w", err)
	}

	err = fs.WalkDir(d.buildFS, ".", func(path string, e fs.DirEntry, _ error) error {
		if !isLevelDeep(path) {
			return nil
		}

		if path == "." {
			return nil
		}

		if e == nil {
			return fmt.Errorf("failed to get dir entry: %w", err)
		}

		if e.IsDir() {
			return nil
		}

		file, err := d.buildFS.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		stat, err := file.Stat()
		if err != nil {
			return fmt.Errorf("failed to stat file: %w", err)
		}

		size := stat.Size()

		bytes := make([]byte, size)

		_, err = file.Read(bytes)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		err = os.WriteFile(d.dir+"/build/"+path, bytes, 0o644)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk build fs: %w", err)
	}

	return nil
}

func isLevelDeep(path string) bool {
	dir, _ := filepath.Split(path)
	return dir == ""
}

func (d *Digraph) Dir() string {
	return d.dir
}
