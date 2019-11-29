package hoistyml

import (
	"github.com/JosiahWitt/erk"
	"github.com/JosiahWitt/erk/erg"
	"gopkg.in/yaml.v3"
)

type (
	ErkInvalidVersion     erk.DefaultKind
	ErkUnsupportedVersion erk.DefaultKind
	ErkInvalidServices    erk.DefaultKind
	ErkInvalidService     erk.DefaultKind
)

var (
	ErrVersionRequired    = erk.New(ErkInvalidVersion{}, "hoist.yml version is required")
	ErrVersionUnsupported = erk.New(ErkUnsupportedVersion{}, "only hoist.yml version 0.1.0 is supported, got '{{.version}}'")
	ErrServiceMissingType = erk.New(ErkInvalidService{}, "service '{{.name}}' missing type")
	ErrServiceMissingPath = erk.New(ErkInvalidService{}, "service '{{.name}}' missing path (use '.' for same directory)")
)

// Template for the hoist.yml format.
type Template struct {
	Version  string   `yaml:"version"`
	Services Services `yaml:"services"`
}

// Services in the Hoist stack.
type Services map[string]*Service

// Service in the Hoist stack.
type Service struct {
	Name string `yaml:"-"` // Added by Services.Parse()
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

// Parse the template and check for errors.
func (t *Template) Parse() error {
	if t.Version == "" {
		return ErrVersionRequired
	} else if t.Version != "0.1.0" {
		return erk.WithParam(ErrVersionUnsupported, "version", t.Version)
	}

	return t.Services.Parse()
}

// String representation of the template in YAML.
func (t *Template) String() string {
	bytes, err := yaml.Marshal(t)
	if err != nil {
		return "could not YAML marshal template: " + err.Error()
	}

	return string(bytes)
}

// Parse the services and check for errors.
func (s Services) Parse() error {
	errs := erg.New(ErkInvalidServices{}, "hoist.yml services invalid")

	for name, service := range s {
		service.Name = name
		if err := service.Parse(); err != nil {
			errs = erg.Append(errs, err)
		}
	}

	if erg.Any(errs) {
		return errs
	}

	return nil
}

// Parse the service and check for errors.
func (s *Service) Parse() error {
	switch {
	case s.Type == "":
		return erk.WithParam(ErrServiceMissingType, "name", s.Name)
	case s.Path == "":
		return erk.WithParam(ErrServiceMissingPath, "name", s.Name)
	default:
		return nil
	}
}
