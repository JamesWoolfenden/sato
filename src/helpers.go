package sato

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/gobeam/stringy"
)

func replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

func add(s string, a []string, m map[string]bool) ([]string, map[string]bool) {
	if m[s] {
		return a, m
	}
	a = append(a, s)
	m[s] = true
	return a, m
}

func split(source string, separator string) []string {
	return strings.Split(source, separator)
}

func dequote(target string) string {
	return strings.Replace(target, "\"", "", -1)
}

func quote(target string) string {
	return "\"" + target + "\""
}

func boolean(test *bool) bool {
	if test == nil {
		return false
	}
	return *test
}

func decode64(str string) string {
	temp, _ := base64.StdEncoding.DecodeString(str)
	return string(temp)
}

func sprint(unknown interface{}) string {
	temp := strings.Replace(strings.Replace(fmt.Sprint(unknown), "[", "", 1), "]", "", 1)
	if temp == "<nil>" {
		return "\"\""
	}
	return temp
}

func snake(Camel string) string {
	str := stringy.New(Camel)
	snakeStr := str.SnakeCase()
	return snakeStr.ToLower()
}

func kebab(Camel string) string {
	str := stringy.New(Camel)
	KebabStr := str.KebabCase()
	return KebabStr.ToLower()
}

func lower(target string) string {
	return strings.ToLower(target)
}

func nill(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func nild(str *string, myDefault string) string {
	if str == nil || *str == "" {
		return myDefault
	}
	return *str
}

func array(mySlice []string) string {
	if mySlice == nil || mySlice[0] == "" {
		return "[]"
	}
	var newSlice []string
	for _, item := range mySlice {
		newSlice = append(newSlice, "\t\""+item+"\"")
	}
	newString := "[\n" + strings.Join(newSlice[:], ",") + "\n\t]\n"
	return newString
}

func arrayReplace(mySlice []string, target string, replacement string) string {
	if mySlice == nil || mySlice[0] == "" {
		return "[]"
	}
	var newSlice []string

	for _, item := range mySlice {
		item = strings.Replace(item, target, replacement, 1)
		newSlice = append(newSlice, "\t\""+item+"\"")
	}
	newString := "[\n" + strings.Join(newSlice[:], ",") + "\n\t]\n"
	return newString
}

func contains(target string, substring string) bool {
	return strings.Contains(target, substring)
}

func zipfile(code string, filename string, runtime string) string {
	var extension string
	switch runtime {
	case "nodejs16.x", "nodejs14.x", "nodejs12.x", "nodejs":
		extension = ".js"
	case "python3.9", "python3.8", "python3.7", "python3.6":
		extension = ".py"
	case "go1.x":
		extension = ".go"
	default:
		extension = ".txt"
	}

	codeFile := filename + extension
	d1 := []byte(code)
	_ = os.WriteFile(codeFile, d1, 0644)

	output := filename + ".zip"
	archive, _ := os.Create(output)
	defer func(archive *os.File) {
		_ = archive.Close()
	}(archive)
	zipWriter := zip.NewWriter(archive)

	f1, _ := os.Open(filename)

	defer func(f1 *os.File) {
		_ = f1.Close()
	}(f1)

	w1, _ := zipWriter.Create(filename)
	_, _ = io.Copy(w1, f1)
	_ = zipWriter.Close()

	return output
}
