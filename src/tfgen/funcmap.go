package tfgen

import (
	"maps"
	"strings"
	"text/template"
)

// Deref returns the value of a string pointer, or empty string if nil.
func Deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// FuncMap returns the base set of template functions used by all converters.
// Any extra maps are merged on top, allowing each converter to register its
// provider-specific helpers (Tags, TFJoin, Uuid, etc) without redeclaring the
// shared core. A fresh map is returned on every call.
func FuncMap(extra ...template.FuncMap) template.FuncMap {
	m := template.FuncMap{
		"Array":        Array,
		"ArrayReplace": ArrayReplace,
		"Boolean":      Boolean,
		"Contains":     Contains,
		"Decode64":     Decode64,
		"Demap":        Demap,
		"Dequote":      Dequote,
		"Deref":        Deref,
		"Kebab":        Kebab,
		"Map":          Map,
		"Marshal":      Marshal,
		"Nil":          Nill,
		"Nild":         Nild,
		"Quote":        Quote,
		"RandomString": RandomString,
		"Replace":      Replace,
		"Snake":        Snake,
		"Split":        Split,
		"SplitOn":      SplitOn,
		"Sprint":       Sprint,
		"ToLower":      Lower,
		"ToUpper":      strings.ToUpper,
		"ZipFile":      Zipfile,
	}
	for _, e := range extra {
		maps.Copy(m, e)
	}
	return m
}
