package see

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func extract(filename string, destination string) {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(destination, f.Name)
		//fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}

		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			_ = os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}

func UpdateSchema(directory string) {
	zip := "schema.zip"
	destination := filepath.Join(directory, zip)
	err := DownloadFile(
		"https://schema.cloudformation.us-east-1.amazonaws.com/CloudformationSchema.zip", destination)

	if err != nil {
		log.Fatal().Msg("failed to update schema")
	}

	extract(destination, directory)
}

func TestLookupAll(t *testing.T) {
	t.Parallel()
	directory := "../../schema"

	UpdateSchema(directory)

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

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
