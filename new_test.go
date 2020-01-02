package rightmpl

import (
	"bytes"
	"io"
	"testing"
)

var testComponents = map[string]func(w io.Writer, values map[string]interface{}){
	"sayBoo": func(w io.Writer, values map[string]interface{}) {
		w.Write([]byte("boo"))
	},
}

func Test_New(t *testing.T) {
	tests := []struct {
		name   string
		tmpl   string
		m      map[string]interface{}
		expect string
	}{{
		"hi",
		"hello.tmpl",
		map[string]interface{}{"val": "world"},
		"hello world",
	}, {
		"addfunc",
		"add.tmpl",
		nil,
		"sub says boo: boo",
	}}
	ts, err := New("testtmpl/", testComponents)
	if err != nil {
		t.Fatal("got: ", err.Error(), " want err=nil ")
	}

	for r := range tests {
		w := bytes.NewBuffer(nil)
		err = ts.Render(w, r.tmpl, r.m)
		if err != nil {
			t.Fatal("got: ", err.Error(), " want err=nil ")
		}
		if w.String() != r.expect {
			t.Error("got: '" + w.String() + "', expected '" + r.expect + "'")
		}
	}
}
