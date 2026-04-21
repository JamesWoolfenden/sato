package tfgen

import "testing"

//goland:noinspection GoLinter
func TestWrite(t *testing.T) {
	t.Parallel()

	type args struct {
		output   string
		location string
		name     string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Pass", args{"gibberine", ".", "test"}, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := Write(tt.args.output, tt.args.location, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
