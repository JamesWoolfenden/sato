package tfgen

import (
	"archive/zip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobeam/stringy"
	"github.com/rs/zerolog/log"
)

// Replace is a template function.
func Replace(input, from, to string) string {
	return strings.ReplaceAll(input, from, to)
}

// Add appends a string to a slice if not already tracked in the seen map.
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
	temp, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Warn().Err(err).Msgf("failed to decode base64 string: %s", str)
		return ""
	}

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

	newSlice := make([]string, 0, len(mySlice))

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

	newSlice := make([]string, 0, len(mySlice))

	for _, item := range mySlice {
		item = strings.Replace(item, target, replacement, 1)
		newSlice = append(newSlice, "\t\""+item+"\"")
	}

	return "[\n" + strings.Join(newSlice, ",") + "\n\t]\n"
}

// Contains is a template function.
func Contains(target string, substring string) bool {
	return strings.Contains(target, substring)
}

// Zipfile creates a zip archive from inline code for use in Terraform lambda/function resources.
func Zipfile(code string, filename string, runtime string, destination ...string) string {
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

	var destDir string
	if len(destination) > 0 && destination[0] != "" {
		destDir = destination[0]
	} else {
		destDir = "."
	}

	codeFileName := filename + extension
	codeFile := filepath.Join(destDir, codeFileName)
	zipFileName := filename + ".zip"
	output := filepath.Join(destDir, zipFileName)

	d1 := []byte(code)
	if err := os.WriteFile(codeFile, d1, 0o600); err != nil {
		log.Error().Err(err).Msgf("failed to write code file: %s", codeFile)
		return ""
	}

	archive, err := os.Create(output) // #nosec G304 -- Creating zip from controlled filename
	if err != nil {
		log.Error().Err(err).Msgf("failed to create zip archive: %s", output)
		return ""
	}

	defer func(archive *os.File) {
		if err := archive.Close(); err != nil {
			log.Warn().Err(err).Msg("failed to close archive")
		}
	}(archive)

	zipWriter := zip.NewWriter(archive)

	file, err := os.Open(codeFile) // #nosec G304 -- Opening file created by this function
	if err != nil {
		log.Error().Err(err).Msgf("failed to open file: %s", codeFile)
		return ""
	}

	defer func(f1 *os.File) {
		if err := f1.Close(); err != nil {
			log.Warn().Err(err).Msg("failed to close file")
		}
	}(file)

	w1, err := zipWriter.Create(codeFileName)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create zip entry: %s", filename)
		return ""
	}

	if _, err := io.Copy(w1, file); err != nil {
		log.Error().Err(err).Msg("failed to copy file to zip")
		return ""
	}

	if err := zipWriter.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close zip writer")
		return ""
	}

	return zipFileName
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

// RandomString is a template function.
func RandomString(n int) string {
	rand.New(rand.NewSource(time.Now().UnixNano())) // #nosec G404 -- Non-crypto random OK for resource naming

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	myString := make([]rune, n)

	for i := range myString {
		myString[i] = letters[rand.Intn(len(letters))] // #nosec G404 -- Non-crypto random OK for resource naming
	}

	return string(myString)
}

// Map is a template function.
func Map(myMap map[string]string) string {
	var builder strings.Builder
	builder.WriteString("{ \n")

	for item, stuff := range myMap {
		builder.WriteString("\t\"")
		builder.WriteString(item)
		builder.WriteString("\"=\"")
		builder.WriteString(stuff)
		builder.WriteString("\"\n")
	}

	builder.WriteString(" }")

	return builder.String()
}

// Marshal is a template function.
//
//goland:noinspection GoUnnecessarilyExportedIdentifiers
func Marshal(v interface{}) string {
	a, err := json.Marshal(v)
	if err != nil {
		log.Error().Msgf("marshalling failed %s", err)
	}

	return string(a)
}
