package see

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestLookup(t *testing.T) {
	t.Parallel()

	type args struct {
		resource string
		reverse  bool
	}

	result := "aws_appautoscaling_target"
	myServiceBus := "azurerm_servicebus_namespace"
	reverse := "aws::efs::filesystem"
	none := "none"

	tests := []struct {
		name    string
		args    args
		want    *string
		wantErr bool
	}{
		{name: "Pass", args: args{"AWS::ApplicationAutoScaling::ScalableTarget", false}, want: &result, wantErr: false},
		{name: "None", args: args{"alexa::ask::skill", false}, want: &none, wantErr: false},
		{name: "Pass", args: args{"Microsoft.ServiceBus/namespaces/", false}, want: &myServiceBus, wantErr: false},
		{name: "Fail", args: args{"AWS::Guff::Guffing", false}, want: nil, wantErr: true},
		{name: "Reverse", args: args{resource: "aws_efs_file_system", reverse: true}, want: &reverse, wantErr: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Lookup(tt.args.resource, tt.args.reverse)

			if (err != nil) != tt.wantErr {
				t.Errorf("Lookup() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lookup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLookupAll(t *testing.T) {
	t.Parallel()
	directory := "../../schema"
	files, err := os.ReadDir(directory)

	if err != nil {
		log.Warn().Msgf("failed to find files at  %s", directory)
	}

	for _, file := range files {
		//has file extension json
		//has lookup
		if strings.Contains(file.Name(), ".json") {
			fileName := filepath.Join(directory, file.Name())
			content, err := os.ReadFile(fileName)
			if err != nil {
				log.Warn().Msgf("file is empty %s", file.Name())
			}

			var result map[string]interface{}
			err = json.Unmarshal(content, &result)

			if err != nil {
				log.Error().Msgf("failed to parse %s", fileName)
				return
			}

			typeName := strings.ToLower(result["typeName"].(string))
			_, err = Lookup(typeName, false)
			if err != nil {

				var s strings.Builder
				s.WriteString("\"")
				s.WriteString(typeName)
				s.WriteString("\": \"\",")
				fmt.Println(s.String())
				t.Errorf("Lookup incomplete")
			}
		}
	}

	//	got, err := Lookup(tt.args.resource, tt.args.reverse)

}
