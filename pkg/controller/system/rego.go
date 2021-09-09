package system

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/pkg/errors"

	"github.com/crossplane-contrib/provider-styra/apis/system/v1alpha1"
)

const (
	errCannotParseLabels         = "cannot parse labels rego"
	errCannotParseLabelsTemplate = "cannot parse labels template"
	errCannotGenerateLabelsRego  = "cannot generate labels.rego"
)

const labelsRegoTemplate = `
package metadata.{{ .Name }}.labels

labels := {
{{ range $k, $v := .Labels -}}
{{ printf "\"%s\": \"%v\"," $k $v | indent 4 }}
{{ end -}}
}
`

func generateRegoLabels(cr *v1alpha1.System) (string, error) {
	temp, err := template.New("labels.rego").Funcs(sprig.TxtFuncMap()).Parse(labelsRegoTemplate)
	if err != nil {
		return "", errors.Wrap(err, errCannotParseLabelsTemplate)
	}

	type templateData struct {
		Name   string
		Labels map[string]string
	}

	labels := make(map[string]string, len(cr.Spec.ForProvider.Labels)+1)
	labels["system-type"] = cr.Spec.ForProvider.Type
	for key, value := range cr.Spec.ForProvider.Labels {
		labels[key] = value
	}

	data := &templateData{
		Name:   meta.GetExternalName(cr),
		Labels: labels,
	}

	buffer := &bytes.Buffer{}
	if err := temp.Execute(buffer, data); err != nil {
		return "", errors.Wrap(err, errCannotGenerateLabelsRego)
	}

	return buffer.String(), nil
}

func compareLabels(ctx context.Context, labelsModule string, cr *v1alpha1.System) (bool, error) {
	parsedLabels, err := extractRegoLabels(ctx, meta.GetExternalName(cr), labelsModule)
	if err != nil {
		return false, errors.Wrap(err, errCannotParseLabels)
	}

	remaining := make(map[string]string, len(cr.Spec.ForProvider.Labels)+1)
	remaining["system-type"] = cr.Spec.ForProvider.Type
	for key, value := range cr.Spec.ForProvider.Labels {
		remaining[key] = value
	}

	for key, value := range parsedLabels {
		specValue, exists := remaining[key]
		if !exists {
			return false, nil
		}

		if specValue != value {
			return false, nil
		}

		delete(remaining, key)
	}

	if len(remaining) > 0 {
		return false, nil
	}

	return true, nil
}

func extractRegoLabels(ctx context.Context, systemID, labelsModule string) (map[string]string, error) {
	// Compile the module. The keys are used as identifiers in error messages.
	compiler, err := ast.CompileModules(map[string]string{
		"labels.rego": labelsModule,
	})

	if err != nil {
		return nil, err
	}

	// // Create a new query that uses the compiled policy from above.
	rego := rego.New(
		rego.Query(fmt.Sprintf("data.metadata.%s.labels.labels", systemID)),
		rego.Compiler(compiler),
	)

	results, err := rego.Eval(ctx)
	if err != nil {
		return nil, err
	}

	labels := make(map[string]string)

	for _, res := range results {
		if len(res.Expressions) == 0 {
			continue
		}

		expression, isMap := res.Expressions[0].Value.(map[string]interface{})
		if !isMap {
			continue
		}

		for labelKey, value := range expression {
			labelValue, isString := value.(string)
			if !isString {
				continue
			}

			labels[labelKey] = labelValue
		}
	}

	return labels, nil
}
