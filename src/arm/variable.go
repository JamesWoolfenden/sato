package arm

import (
	_ "embed" //required for embed
)

//go:embed variable.template
var VariableFile []byte
