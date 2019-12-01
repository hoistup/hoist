package hoistyml_test

import (
	"fmt"
	"testing"

	"github.com/hoistup/hoist/cmd/hoist/hoistyml"
	"github.com/matryer/is"
)

func TestTemplateString(t *testing.T) {
	t.Run("with valid template", func(t *testing.T) {
		is := is.New(t)

		tmpl := &hoistyml.Template{
			Version: "0.1.0",
			Stack: hoistyml.Stack{
				Name: "my-stack",
			},
			Services: hoistyml.Services{
				"svc1": {
					Name: "svc1",
					Type: "go",
					Path: "my/svc1",
				},
			},
		}

		is.Equal(fmt.Sprintf("%v", tmpl),
			`version: 0.1.0
stack:
    name: my-stack
services:
    svc1:
        type: go
        path: my/svc1
`,
		)
	})
}

// Note: The other methods are tested in load_test.go
