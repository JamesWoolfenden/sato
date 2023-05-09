package arm

import (
	_ "embed" //required for embed
)

//go:embed variable.template
var VariableFile []byte

//go:embed output.template
var OutputFile []byte
