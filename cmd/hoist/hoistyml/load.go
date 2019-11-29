package hoistyml

import (
	"path/filepath"

	"github.com/JosiahWitt/erk"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type (
	ErkUnableToReadFile  erk.DefaultKind
	ErkUnableToUnmarshal erk.DefaultKind
)

var (
	ErrUnableToReadFile  = erk.New(ErkUnableToReadFile{}, "cannot read '{{.path}}'; make sure you are in the correct directory")
	ErrUnableToUnmarshal = erk.New(ErkUnableToReadFile{}, "cannot unmarshal '{{.path}}' into the hoist.yml format: {{.err}}")
)

// FileName for the hoist.yml file.
const FileName = "hoist.yml"

// Load the hoist.yml file from the provided file system in the given directory.
func Load(fs afero.Fs, dir string) (*Template, error) {
	path := filepath.Join(dir, FileName)
	bytes, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, erk.WrapAs(erk.WithParam(ErrUnableToReadFile, "path", path), err)
	}

	var tmpl Template
	if err := yaml.Unmarshal(bytes, &tmpl); err != nil {
		return nil, erk.WrapAs(erk.WithParam(ErrUnableToUnmarshal, "path", path), err)
	}

	if err := tmpl.Parse(); err != nil {
		return nil, err
	}

	return &tmpl, nil
}
