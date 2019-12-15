package hoistyml

import (
	"regexp"

	"github.com/JosiahWitt/erk"
	"github.com/JosiahWitt/erk/erg"
	"github.com/hoistup/hoist/cmd/hoist/erks"
	"gopkg.in/yaml.v3"
)

type (
	ErkInvalidVersion     struct{ erks.Default }
	ErkUnsupportedVersion struct{ erks.Default }
	ErkInvalidStack       struct{ erks.Default }
	ErkInvalidServices    struct{ erks.Default }
	ErkInvalidService     struct{ erks.Default }
)

var (
	ErrVersionRequired    = erk.New(ErkInvalidVersion{}, "hoist.yml 'version' is required")
	ErrVersionUnsupported = erk.New(ErkUnsupportedVersion{}, "only hoist.yml version 0.1.0 is supported, got '{{.version}}'")

	ErrStackMissingName = erk.New(ErkInvalidStack{}, "hoist.yml 'stack' missing 'name'")
	ErrStackNameInvalid = erk.New(ErkInvalidStack{}, "hoist.yml stack name, '{{.name}}', invalid; it can only contain a-z, 0-9, and '-'")

	ErrServiceNameInvalid = erk.New(ErkInvalidService{}, "service name '{{.name}}' is invalid; it can only contain a-z, 0-9, and '-'")
	ErrServiceMissingType = erk.New(ErkInvalidService{}, "service '{{.name}}' missing 'type'")
	ErrServiceMissingPath = erk.New(ErkInvalidService{}, "service '{{.name}}' missing 'path' (use '.' for same directory)")
)

// Regexp for names, since they should be valid subdomains
var nameRegexp = regexp.MustCompile("^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$")

// Template for the hoist.yml format.
type Template struct {
	Version  string   `yaml:"version"`
	Stack    Stack    `yaml:"stack"`
	Services Services `yaml:"services"`
}

// Stack details.
type Stack struct {
	Name string `yaml:"name"`
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

	if err := t.Stack.Parse(); err != nil {
		return err
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

// Parse the stack and check for errors.
func (s *Stack) Parse() error {
	if s.Name == "" {
		return ErrStackMissingName
	}

	// Stack names should be valid subdomains
	if !nameRegexp.MatchString(s.Name) {
		return erk.WithParam(ErrStackNameInvalid, "name", s.Name)
	}

	return nil
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
	case !nameRegexp.MatchString(s.Name):
		return erk.WithParam(ErrServiceNameInvalid, "name", s.Name)
	case s.Type == "":
		return erk.WithParam(ErrServiceMissingType, "name", s.Name)
	case s.Path == "":
		return erk.WithParam(ErrServiceMissingPath, "name", s.Name)
	default:
		return nil
	}
}
