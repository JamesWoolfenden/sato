package tfgen

import (
	"testing"
	"text/template"
)

func TestFuncMap(t *testing.T) {
	t.Parallel()

	base := FuncMap()
	for _, k := range []string{"Array", "Dequote", "Marshal", "Snake", "ZipFile", "Deref"} {
		if _, ok := base[k]; !ok {
			t.Errorf("FuncMap() missing base key %q", k)
		}
	}

	extra := FuncMap(template.FuncMap{"Tags": func() string { return "" }})
	if _, ok := extra["Tags"]; !ok {
		t.Error("FuncMap(extra) did not merge extra key")
	}
	if _, ok := base["Tags"]; ok {
		t.Error("FuncMap() leaked extra key into earlier result")
	}
}

func TestDeref(t *testing.T) {
	t.Parallel()

	s := "x"
	if got := Deref(&s); got != "x" {
		t.Errorf("Deref(&s) = %q, want %q", got, "x")
	}
	if got := Deref(nil); got != "" {
		t.Errorf("Deref(nil) = %q, want empty", got)
	}
}
