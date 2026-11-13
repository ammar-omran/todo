package render_test

import (
	"fmt"
	"html/template"
	"strings"
	"sync"
	"testing"
	"testing/fstest"

	"to-do/internal/render"
)

func TestRenderReplacesTFuncPerCall(t *testing.T) {
	fsys := fstest.MapFS{
		"page.html": &fstest.MapFile{
			Data: []byte(`<h1>{{ T "greeting" }}</h1><p>{{ T "count" 3 }}</p>`),
		},
	}

	renderer := render.NewRegistry().LoadFS(fsys, "page.html")

	first, err := renderer.Render(nil, template.FuncMap{
		"T": func(key string, args ...any) string {
			return "EN:" + key + fmt.Sprint(args)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := "<h1>EN:greeting[]</h1><p>EN:count[3]</p>"; first != want {
		t.Fatalf("first render = %q, want %q", first, want)
	}

	second, err := renderer.Render(nil, template.FuncMap{
		"T": func(key string, args ...any) string {
			return "AR:" + key + fmt.Sprint(args)
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := "<h1>AR:greeting[]</h1><p>AR:count[3]</p>"; second != want {
		t.Fatalf("second render = %q, want %q", second, want)
	}
}

func TestRenderSharedDefinitionUsesRequestFuncs(t *testing.T) {
	fsys := fstest.MapFS{
		"list.html": &fstest.MapFile{
			Data: []byte(`<ul>{{ range . }}<li>{{ T "item" . }}</li>{{ end }}</ul>`),
		},
	}

	renderer := render.NewRegistry().LoadFS(fsys, "list.html")

	html, err := renderer.Render([]string{"a", "b"}, template.FuncMap{
		"T": func(key string, args ...any) string {
			return fmt.Sprintf("%s:%v", key, args[0])
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := "<ul><li>item:a</li><li>item:b</li></ul>"; html != want {
		t.Fatalf("render = %q, want %q", html, want)
	}
}

func TestRenderConcurrentCalls(t *testing.T) {
	fsys := fstest.MapFS{
		"page.html": &fstest.MapFile{
			Data: []byte(`{{ T "value" }}`),
		},
	}

	renderer := render.NewRegistry().LoadFS(fsys, "page.html")

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			html, err := renderer.Render(nil, template.FuncMap{
				"T": func(string, ...any) string {
					return fmt.Sprintf("value-%d", i)
				},
			})
			if err != nil {
				t.Errorf("render %d: %v", i, err)
				return
			}
			if !strings.HasPrefix(html, "value-") {
				t.Errorf("render %d returned %q", i, html)
			}
		}(i)
	}
	wg.Wait()
}
