package sato

import "testing"

func Test_replace(t *testing.T) {
	type args struct {
		input string
		from  string
		to    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replace(tt.args.input, tt.args.from, tt.args.to); got != tt.want {
				t.Errorf("replace() = %v, want %v", got, tt.want)
			}
		})
	}
}
