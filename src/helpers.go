package sato

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/awslabs/goformation/v7/cloudformation/tags"
	"github.com/gobeam/stringy"
)

// Replace is a template function
func Replace(input, from, to string) string {
	return strings.Replace(input, from, to, -1)
}

// Add is a template function
func Add(s string, a []string, m map[string]bool) ([]string, map[string]bool) {
	if m[s] {
		return a, m
	}
	a = append(a, s)
	m[s] = true
	return a, m
}

// Split is a template function
func Split(source string, separator string) []string {
	return strings.Split(source, separator)
}

// Dequote is a template function
func Dequote(target string) string {
	return strings.Replace(target, "\"", "", -1)
}

// Quote is a template function
func Quote(target string) string {
	return "\"" + target + "\""
}

// Boolean is a template function
func Boolean(test *bool) bool {
	if test == nil {
		return false
	}
	return *test
}

// Decode64 is a template function
func Decode64(str string) string {
	temp, _ := base64.StdEncoding.DecodeString(str)
	return string(temp)
}

// Sprint is a template function
func Sprint(unknown interface{}) string {
	temp := strings.Replace(strings.Replace(fmt.Sprint(unknown), "[", "", 1), "]", "", 1)
	if temp == "<nil>" {
		return "\"\""
	}
	return temp
}

// Snake is a template function
func Snake(Camel string) string {
	str := stringy.New(Camel)
	snakeStr := str.SnakeCase()
	return snakeStr.ToLower()
}

// Kebab is a template function
func Kebab(Camel string) string {
	str := stringy.New(Camel)
	KebabStr := str.KebabCase()
	return KebabStr.ToLower()
}

// Lower is a template function
func Lower(target string) string {
	return strings.ToLower(target)
}

// Nill is a template function
func Nill(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

// Nild is a template function
func Nild(str *string, myDefault string) string {
	if str == nil || *str == "" {
		return myDefault
	}
	return *str
}

// Array is a template function
func Array(mySlice []string) string {
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

// ArrayReplace is a template function
func ArrayReplace(mySlice []string, target string, replacement string) string {
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

// Contains is a template function
func Contains(target string, substring string) bool {
	return strings.Contains(target, substring)
}

// Zipfile is a template function
func Zipfile(code string, filename string, runtime string) string {
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

// Demap is a template function
func Demap(str string) []string {
	str = strings.Replace(str, "{", "", -1)
	str = strings.Replace(str, "}", "", -1)
	str = strings.Replace(str, "\"", "", -1)
	str = strings.Replace(str, ":", "", -1)
	str = strings.Replace(str, " ", "", -1)
	return strings.Split(str, ",")
}

// Tags is a template function
func Tags(v []tags.Tag) string {
	var temp string
	for _, item := range v {
		if item.Key != "" {
			temp = temp + "\t\"" + item.Key + "\"" + "=" + "\"" + item.Value + "\"" + "\n"
		}
	}
	return temp
}

// RandomString is a template function
func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Map is a template function
func Map(myMap map[string]string) string {
	result := "{ \n"
	for item, stuff := range myMap {
		result = result + "\t\"" + item + "\"" + "=" + "\"" + stuff + "\"\n"
	}
	result = result + " }"

	return result
}
