package promote

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func write(t *testing.T, path, body string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRun_cf(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	write(t, filepath.Join(root, "src/cf/resources/.keep"), "")
	write(t, filepath.Join(root, "src/cf/resources.go"), "package cf\n\nimport _ \"embed\"\n")
	write(t, filepath.Join(root, "src/cf/resource_mapping.go"),
		"package cf\n\nvar tfLookup = map[string]interface{}{\n\t\"AWS::SNS::Topic\": awsSnsTopic,\n}\n")
	write(t, filepath.Join(root, "src/see/resource_mapping.go"),
		"package see\n\nvar lookupMapping = map[string]string{\n\t\"aws::scheduler::schedule\": none,\n}\n")

	draft := filepath.Join(root, "aws_scheduler_schedule.template")
	write(t, draft, "resource \"aws_scheduler_schedule\" \"{{.item}}\" {}\n")

	if err := Run(draft, "AWS::Scheduler::Schedule", root); err != nil {
		t.Fatalf("Run: %v", err)
	}

	for path, want := range map[string]string{
		"src/cf/resources/aws_scheduler_schedule.template": "{{.item}}",
		"src/cf/resources.go":                              "var awsSchedulerSchedule []byte",
		"src/cf/resource_mapping.go":                       "\"AWS::Scheduler::Schedule\": awsSchedulerSchedule,",
		"src/see/resource_mapping.go":                      "\"aws::scheduler::schedule\": \"aws_scheduler_schedule\",",
	} {
		got, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}

		if !strings.Contains(string(got), want) {
			t.Errorf("%s missing %q, got:\n%s", path, want, got)
		}
	}

	if err := Run(draft, "AWS::Scheduler::Schedule", root); err == nil {
		t.Error("second Run should fail (template already exists)")
	}
}

func Test_camel(t *testing.T) {
	t.Parallel()

	if got := camel("aws_scheduler_schedule"); got != "awsSchedulerSchedule" {
		t.Errorf("camel = %q", got)
	}
}
