package sato

import (
	"encoding/base64"
	"fmt"
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
