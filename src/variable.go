package sato

import (
	_ "embed" //required for embed
)

//go:embed variable.template
var variableFile []byte
