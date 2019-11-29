package hoistyml_test

import (
	"errors"
	"testing"

	"github.com/hoistup/hoist/cmd/hoist/hoistyml"
	"github.com/lithammer/dedent"
	"github.com/matryer/is"
	"github.com/spf13/afero"
)

func TestLoad(t *testing.T) {
	table := []struct {
		Name             string
		Dir              string
		Before           func(fs afero.Fs) error
		ExpectedTemplate *hoistyml.Template
		ExpectedError    error
	}{
		{
			Name: "with valid file",
			Dir:  "my/stack",
			Before: func(fs afero.Fs) error {
				data := dedent.Dedent(`
					version: 0.1.0

					services:
					  myService:
					    type: go
					    path: my/path
				`)

				return afero.WriteFile(fs, "my/stack/hoist.yml", []byte(data), 0655)
			},
			ExpectedTemplate: &hoistyml.Template{
				Version: "0.1.0",
				Services: hoistyml.Services{
					"myService": {
						Name: "myService",
						Type: "go",
						Path: "my/path",
					},
				},
			},
		},
		{
			Name:          "with missing file",
			Dir:           "my/stack",
			ExpectedError: hoistyml.ErrUnableToReadFile,
		},
		{
			Name: "with invalid YAML",
			Dir:  "my/invalid-stack",
			Before: func(fs afero.Fs) error {
				data := "!!!THIS ISN'T YAML!!!"

				return afero.WriteFile(fs, "my/invalid-stack/hoist.yml", []byte(data), 0655)
			},
			ExpectedError: hoistyml.ErrUnableToUnmarshal,
		},
		{
			Name: "with missing version",
			Dir:  "my/invalid-stack",
			Before: func(fs afero.Fs) error {
				data := dedent.Dedent(`
					services:
					  myService:
					    type: go
					    path: my/path
				`)

				return afero.WriteFile(fs, "my/invalid-stack/hoist.yml", []byte(data), 0655)
			},
			ExpectedError: hoistyml.ErrVersionRequired,
		},
		{
			Name: "with unsupported version",
			Dir:  "my/invalid-stack",
			Before: func(fs afero.Fs) error {
				data := dedent.Dedent(`
					version: 0.1.1

					services:
					  myService:
					    type: go
					    path: my/path
				`)

				return afero.WriteFile(fs, "my/invalid-stack/hoist.yml", []byte(data), 0655)
			},
			ExpectedError: hoistyml.ErrVersionUnsupported,
		},
		{
			Name: "with invalid service: missing type",
			Dir:  "my/invalid-stack",
			Before: func(fs afero.Fs) error {
				data := dedent.Dedent(`
					version: 0.1.0

					services:
					  myService:
					    path: my/path
				`)

				return afero.WriteFile(fs, "my/invalid-stack/hoist.yml", []byte(data), 0655)
			},
			ExpectedError: hoistyml.ErrServiceMissingType,
		},
		{
			Name: "with invalid service: missing path",
			Dir:  "my/invalid-stack",
			Before: func(fs afero.Fs) error {
				data := dedent.Dedent(`
					version: 0.1.0

					services:
					  myService:
					    type: go
				`)

				return afero.WriteFile(fs, "my/invalid-stack/hoist.yml", []byte(data), 0655)
			},
			ExpectedError: hoistyml.ErrServiceMissingPath,
		},
	}

	for _, entry := range table {
		t.Run(entry.Name, func(t *testing.T) {
			is := is.New(t)
			fs := afero.NewMemMapFs()

			if entry.Before != nil {
				is.NoErr(entry.Before(fs))
			}

			tmpl, err := hoistyml.Load(fs, entry.Dir)
			is.True(errors.Is(err, entry.ExpectedError))
			is.Equal(tmpl, entry.ExpectedTemplate)
		})
	}
}
