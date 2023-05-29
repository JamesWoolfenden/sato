package arm

import (
	_ "embed" // required for embed
)

//go:embed variable.template
var variableFile []byte

//go:embed output.template
var outputFile []byte

//go:embed data.template
var dataFile []byte
