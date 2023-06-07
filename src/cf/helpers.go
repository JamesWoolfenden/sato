package cf

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

// Replace is a template function.
func Replace(input, from, to string) string {
	return strings.ReplaceAll(input, from, to)
}

// Add is a template function.
func Add(
	myString string,
	myList []string,
	myMap map[string]bool,
) ([]string, map[string]bool) {
	if myMap[myString] {
		return myList, myMap
	}
	myList = append(myList, myString) //nolint:wsl
	myMap[myString] = true

	return myList, myMap
}

// Split is a template function.
func Split(source string, separator string) []string {
	return strings.Split(source, separator)
}

// SplitOn is a template function.
func SplitOn(source string, separator string, index int) string {
	splits := strings.Split(source, separator)
	if len(splits) >= index+1 {
		return splits[index]
	}

	return ""
}

// Dequote is a template function.
func Dequote(target string) string {
	return strings.ReplaceAll(strings.ReplaceAll(target, "\"", ""), "'", "")
}

// Quote is a template function.
func Quote(target string) string {
	// is it a resource or variable
	if (strings.Contains(target, "var.") || strings.Contains(target, "local.") ||
		(strings.Contains(target, "_") && strings.Contains(target, "."))) && (!strings.Contains(target, "${")) {
		return target
	}

	return "\"" + target + "\""
}

// Boolean is a template function.
func Boolean(test *bool) bool {
	if test == nil {
		return false
	}

	return *test
}

// Decode64 is a template function.
func Decode64(str string) string {
	temp, _ := base64.StdEncoding.DecodeString(str)

	return string(temp)
}

// Sprint is a template function.
func Sprint(unknown interface{}) string {
	temp := strings.Replace(strings.Replace(fmt.Sprint(unknown), "[", "", 1), "]", "", 1)
	if temp == "<nil>" {
		return "\"\""
	}

	return temp
}

// Snake is a template function.
func Snake(camel string) string {
	str := stringy.New(camel)
	snakeStr := str.SnakeCase()

	return snakeStr.ToLower()
}

// Kebab is a template function.
func Kebab(camel string) string {
	str := stringy.New(camel)
	KebabStr := str.KebabCase()

	return KebabStr.ToLower()
}

// Lower is a template function.
func Lower(target string) string {
	return strings.ToLower(target)
}

// Nill is a template function.
func Nill(str *string) string {
	if str == nil {
		return ""
	}

	return *str
}

// Nild is a template function.
func Nild(str *string, myDefault string) string {
	if str == nil || *str == "" {
		return myDefault
	}

	return *str
}

// Array is a template function.
func Array(mySlice []string) string {
	if mySlice == nil || mySlice[0] == "" {
		return "[]"
	}

	var newSlice []string //nolint:prealloc

	for _, item := range mySlice {
		newSlice = append(newSlice, "\t\""+item+"\"")
	}

	newString := "[\n" + strings.Join(newSlice, ",") + "\n\t]\n"

	return newString
}

// ArrayReplace is a template function.
func ArrayReplace(mySlice []string, target string, replacement string) string {
	if mySlice == nil || mySlice[0] == "" {
		return "[]"
	}

	var newSlice []string //nolint:prealloc

	for _, item := range mySlice {
		item = strings.Replace(item, target, replacement, 1)
		newSlice = append(newSlice, "\t\""+item+"\"")
	}
	newString := "[\n" + strings.Join(newSlice, ",") + "\n\t]\n"

	return newString
}

// Contains is a template function.
func Contains(target string, substring string) bool {
	return strings.Contains(target, substring)
}

// Zipfile is a template function.
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
	_ = os.WriteFile(codeFile, d1, 0o644)

	output := filename + ".zip"
	archive, _ := os.Create(output)

	defer func(archive *os.File) {
		_ = archive.Close()
	}(archive)

	zipWriter := zip.NewWriter(archive)

	file, _ := os.Open(filename)

	defer func(f1 *os.File) {
		_ = f1.Close()
	}(file)

	w1, _ := zipWriter.Create(filename)
	_, _ = io.Copy(w1, file)
	_ = zipWriter.Close()

	return output
}

// Demap is a template function.
func Demap(str string) []string {
	str = strings.ReplaceAll(str, "{", "")
	str = strings.ReplaceAll(str, "}", "")
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, ":", "")
	str = strings.ReplaceAll(str, " ", "")

	return strings.Split(str, ",")
}

// Tags is a template function.
func Tags(v []tags.Tag) string {
	var temp string

	for _, item := range v {
		if item.Key != "" {
			temp = temp + "\t\"" + item.Key + "\"" + "=" + "\"" + item.Value + "\"" + "\n"
		}
	}

	return temp
}

// RandomString is a template function.
func RandomString(n int) string {
	rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	myString := make([]rune, n)

	for i := range myString {
		//goland:noinspection GoLinter
		myString[i] = letters[rand.Intn(len(letters))] //nolint:gosec
	}

	return string(myString)
}

// Map is a template function.
func Map(myMap map[string]string) string {
	result := "{ \n"
	for item, stuff := range myMap {
		result = result + "\t\"" + item + "\"" + "=" + "\"" + stuff + "\"\n"
	}

	result += " }"

	return result
}
