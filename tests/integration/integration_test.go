package integration_test

import (
	"os"
	"path/filepath"
	"testing"

	"sato/src/arm"
	"sato/src/cf"
	"sato/tests/utils"
)

type fixture struct {
	name     string
	source   string
	parse    func(file, dest string) error
	validate bool   // also run `tofu validate` (semantic check)
	skip     string // non-empty: known gap, skip with this reason
}

var fixtures = []fixture{
	{name: "cf-template", source: "../../examples/template.yaml", parse: cf.Parse, validate: true},
	{name: "cf-athena", source: "../../examples/athena.yaml", parse: cf.Parse, validate: true},
	{name: "cf-kinesis", source: "../../examples/kinesis.yaml", parse: cf.Parse, validate: true},
	{name: "cf-aws-vpc", source: "../../examples/aws-vpc.template.yaml", parse: cf.Parse, validate: false},
	{
		name: "cf-linux-bastion", source: "../../examples/linux-bastion-master.template.yaml", parse: cf.Parse,
		skip: "goformation: cannot unmarshal number into Parameters string",
	},
	{
		name:   "arm-vm-simple-windows",
		source: "../../examples/arm/microsoft.compute/vm-simple-windows/azuredeploy.json", parse: arm.Parse,
	},
}

// TestEndToEnd converts every example template and checks the emitted HCL is
// at minimum syntactically valid (tofu init succeeds). Fixtures marked
// validate=true must additionally pass `tofu validate`.
func TestEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping end-to-end tofu tests in -short mode")
	}

	tofu := utils.TofuPath(t)

	for _, f := range fixtures {
		f := f
		t.Run(f.name, func(t *testing.T) {
			if f.skip != "" {
				t.Skip(f.skip)
			}

			dest := t.TempDir()

			if err := f.parse(f.source, dest); err != nil {
				t.Fatalf("parse %s: %v", f.source, err)
			}

			entries, err := os.ReadDir(dest)
			if err != nil {
				t.Fatalf("read %s: %v", dest, err)
			}

			var tfCount int
			for _, e := range entries {
				if filepath.Ext(e.Name()) == ".tf" {
					tfCount++
				}
			}
			if tfCount == 0 {
				t.Fatalf("no .tf files emitted to %s", dest)
			}

			utils.Init(t, tofu, dest)
			if f.validate {
				utils.Validate(t, tofu, dest)
			}
		})
	}
}
