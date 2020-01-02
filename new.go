package rightmpl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/valyala/fasttemplate"
)

// New returns a templates
func New(folderpath string, funcs map[string]Component) (*Templates, error) {
	t := Templates{
		path:  folderpath,
		funcs: map[string]Component{},
		ts:    map[string]*fasttemplate.Template{},
	}
	// FUTURE: add default funcs here. (IF, Template, range)
	for k, v := range funcs {
		t.funcs[k] = v
	}
	err := filepath.Walk(folderpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		return t.addTemplate(path)
	})
	if err != nil {
		return nil, err
	}
	w := watcher.New()
	go func() {
		for {
			select {
			case event := <-w.Event:
				t.tsMx.Lock()
				t.addTemplate(event.Path)
				t.tsMx.Unlock()
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()
	w.AddRecursive(folderpath)
	go func() {
		if err := w.Start(time.Millisecond * 100); err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}()
	return &t, nil
}

// Templates is returned from New to provide access
// to a renderer.
type Templates struct {
	ts    map[string]*fasttemplate.Template
	funcs map[string]Component
	tsMx  sync.RWMutex
	path  string
}

// Component is the signature to add components
type Component func(w io.Writer, values map[string]interface{})

// AddComponent adds a composable function.
func (t *Templates) AddComponent(name string, c Component) {
	t.funcs[name] = c
}

func (t *Templates) addTemplate(path string) error {
	if len(path) < len(t.path) || t.path != path[:len(t.path)] {
		return errors.New("different path")
	}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	// FUTURE: allow richer parse trees by using a parser (antlr?)
	// then turn it into a slice of dynamic function invocations.
	tr, err := fasttemplate.NewTemplate(string(f), "{{", "}}")
	if err != nil {
		return err
	}
	t.ts[path[len(t.path):]] = tr
	return nil
}

var emptyMap = map[string]interface{}{}

// Render a template with values to w.
// pathFrac is the full path after the path given to New().
// values are for print only.
func (t *Templates) Render(w io.Writer, pathFrac string, values map[string]interface{}) error {
	t.tsMx.Lock()
	defer t.tsMx.Unlock()
	if values == nil {
		values = emptyMap
	}
	return t.render(w, pathFrac, values)
}
func (t *Templates) render(w io.Writer, pathFrac string, values map[string]interface{}) error {
	fs, ok := t.ts[pathFrac]
	if !ok {
		return errors.New("Template " + pathFrac + " not found.")
	}
	// TODO do invocations by count instead of tag so we already know what to do .
	fs.ExecuteFunc(w, func(w io.Writer, tag string) (int, error) {
		if v, ok := values[tag]; ok {
			return fmt.Fprint(w, v)
		}
		idx := strings.Index(tag, " ") /// TODO move parsing to the parse step. cache results.
		if idx == -1 {
			idx = len(tag)
		}
		fn, ok := t.funcs[tag[:idx]]
		if !ok {
			return 0, errors.New("Func " + tag[:idx] + " now found.")
		}
		fn(w, values)
		return 0, nil
	})
	return nil
}

/*		argsAry := [4]reflect.Value{wv, nil, nil, nil}
		args := argsAry[1:0]
		for {
			tag = tag[idx:]
			idx = strings.Index(tag, " ")
			if idx == -1 {
				break
			}
			args = append(args, getValue(values, tag[:idx]))
		}
		fn.Call(args)
*/
