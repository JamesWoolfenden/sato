package cf

//goland:noinspection GoLinter
import (
	_ "embed" // required for embed
)

//go:embed variable.template
var variableFile []byte
