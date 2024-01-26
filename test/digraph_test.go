package tests

import (
	"embed"
	"encoding/json"
	"io/fs"
	"os"
	"reflect"
	"sort"
	"testing"

	digraph "github.com/vimcki/go-di-graph"
	"github.com/vimcki/go-di-graph/internal/enhancer"
	"github.com/vimcki/go-di-graph/internal/frontend/d2"
	"github.com/wI2L/jsondiff"
)

type Test struct {
	name       string
	entrypoint string
	fs         embed.FS
	path       string
}

//go:embed _modules/test1/build
var test1 embed.FS

//go:embed _modules/test_set/build
var test2 embed.FS

//go:embed _modules/test_var/build
var test3 embed.FS

//go:embed _modules/nested_structs/build
var test4 embed.FS

//go:embed _modules/map_insert/build
var test5 embed.FS

func TestDigraph(t *testing.T) {
	tests := []Test{
		{
			name:       "test1",
			entrypoint: "Build",
			fs:         test1,
			path:       "_modules/test1/",
		},
		{
			name:       "test_set_build",
			entrypoint: "Set.Build",
			fs:         test2,
			path:       "_modules/test_set/",
		},
		{
			name:       "test_variables",
			entrypoint: "Set.Build",
			fs:         test3,
			path:       "_modules/test_var/",
		},
		{
			name:       "test nested structs",
			entrypoint: "Worker.Fill",
			fs:         test4,
			path:       "_modules/nested_structs/",
		},
		{
			name:       "test inserting into map",
			entrypoint: "Build",
			fs:         test5,
			path:       "_modules/map_insert/",
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			realFS, err := fs.Sub(test.fs, test.path+"build")
			if err != nil {
				t.Fatal(err)
			}

			bytes, err := os.ReadFile(test.path + "config.json")
			if err != nil {
				t.Fatal(err)
			}

			var config interface{}

			err = json.Unmarshal(bytes, &config)
			if err != nil {
				t.Fatal(err)
			}

			dg, err := digraph.New(
				config,
				test.entrypoint,
				realFS,
				digraph.WithNoRender(),
				digraph.WithBlockingHandler(),
			)
			if err != nil {
				t.Fatal(err)
			}

			handler := dg.Handler()

			result, contentType := handler()

			if contentType != "application/json" {
				t.Log(dg.Dir())
				t.Log(string(result))
				t.Fatalf("expected application/json, got %s", contentType)
			}

			var got enhancer.Dependency

			err = json.Unmarshal(result, &got)
			if err != nil {
				t.Fatal(err)
			}

			data, err := os.ReadFile(test.path + "result.json")
			if err != nil {
				t.Fatal(err)
			}

			var want enhancer.Dependency

			err = json.Unmarshal(data, &want)
			if err != nil {
				t.Fatal(err)
			}

			sortTree(&got)
			sortTree(&want)
			if !reflect.DeepEqual(got, want) {
				t.Log(dg.Dir())
				pretty, err := json.MarshalIndent(got, "", "  ")
				if err != nil {
					t.Fatal(err)
				}

				t.Log(string(pretty))

				writeD2(t, pretty)

				gotSorted, err := json.Marshal(got)
				if err != nil {
					t.Fatal(err)
				}

				var gotMap map[string]interface{}

				err = json.Unmarshal(gotSorted, &gotMap)
				if err != nil {
					t.Fatal(err)
				}

				wantSorted, err := json.Marshal(want)
				if err != nil {
					t.Fatal(err)
				}

				var wantMap map[string]interface{}

				err = json.Unmarshal(wantSorted, &wantMap)
				if err != nil {
					t.Fatal(err)
				}

				printDiff(t, wantMap, gotMap)
				t.Fatalf("%v failed", test.name)
			}
		})
	}
}

func printDiff(t *testing.T, want, got map[string]interface{}) {
	patches, err := jsondiff.Compare(got, want)
	if err != nil {
		t.Fatalf("jsondiff error: %v", err)
	}
	for _, patch := range patches {
		t.Logf("%v", patch)
	}
}

func writeD2(t *testing.T, data []byte) {
	data, err := d2.Render(string(data))
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("render.d2", []byte(data), 0o644)
	if err != nil {
		t.Fatal(err)
	}
}

type Dependencies []*enhancer.Dependency

func (d Dependencies) Len() int {
	return len(d)
}

func (d Dependencies) Less(i, j int) bool {
	return d[i].Name < d[j].Name
}

func (d Dependencies) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func sortTree(tree *enhancer.Dependency) {
	if tree == nil {
		return
	}

	for _, dep := range tree.Deps {
		sortTree(dep)
	}

	if tree.Deps != nil {
		sort.Sort(Dependencies(tree.Deps))
	}
}
