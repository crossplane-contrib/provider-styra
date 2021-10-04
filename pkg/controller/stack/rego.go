package stack

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/open-policy-agent/opa/ast"
	"github.com/pkg/errors"

	"github.com/crossplane-contrib/provider-styra/apis/stack/v1alpha1"
)

const (
	errParseSelectors         = "cannot parse selectors rego"
	errParseSelectorsTemplate = "cannot parse selectors template"
	errGenerateSelectorsRego  = "cannot generate selectors.rego"
)

const regoSelectorTemplate = `
package stacks.{{ .Name }}.selectors
import data.library.v1.utils.labels.match.v1 as match

systems[system_id] {
  include := {
{{ range $k, $values := .Include -}}
{{- printf "\"%s\": {" $k | indent 4 }}
{{ range $v := $values }}
{{- printf "\"%s\"," $v | indent 6 }}
{{ end -}}
{{ "}," | indent 4 }}
{{- end }}
  }

  exclude := {
{{ range $k, $values := .Exclude -}}
{{- printf "\"%s\": {" $k | indent 4 }}
{{ range $v := $values }}
{{- printf "\"%s\"," $v | indent 6 }}
{{ end -}}
{{ "}," | indent 4 }}
{{- end }}
  }

  metadata := data.metadata[system_id]
  match.all(metadata.labels.labels, include, exclude)
}
`

type selector map[string][]string

func generateRegoSelectors(cr *v1alpha1.Stack) (string, error) {
	temp, err := template.New("labels.rego").Funcs(sprig.TxtFuncMap()).Parse(regoSelectorTemplate)
	if err != nil {
		return "", errors.Wrap(err, errParseSelectorsTemplate)
	}

	type templateData struct {
		Name    string
		Include selector
		Exclude selector
	}

	data := &templateData{
		Name:    meta.GetExternalName(cr),
		Include: cr.Spec.ForProvider.SelectorInclude,
		Exclude: cr.Spec.ForProvider.SelectorExclude,
	}

	buffer := &bytes.Buffer{}
	if err := temp.Execute(buffer, data); err != nil {
		return "", errors.Wrap(err, errGenerateSelectorsRego)
	}

	return buffer.String(), nil
}

func compareSelectors(selectorsModule string, cr *v1alpha1.Stack) (bool, error) {
	include, exclude, err := extractRegoSelectors(selectorsModule)
	if err != nil {
		return false, errors.Wrap(err, errParseSelectors)
	}

	return isEqualSelector(cr.Spec.ForProvider.SelectorInclude, include) && isEqualSelector(cr.Spec.ForProvider.SelectorExclude, exclude), nil
}

func isEqualSelector(spec, current selector) bool {
	if len(spec) != len(current) {
		return false
	}
	for key, specVal := range spec {
		currentVal, exists := current[key]
		if !exists || len(currentVal) != len(specVal) {
			return false
		}
		specMap := make(map[string]struct{}, len(specVal))
		for _, specLabel := range specVal {
			specMap[specLabel] = struct{}{}
		}
		for _, currentLabel := range currentVal {
			if _, exists := specMap[currentLabel]; !exists {
				return false
			}
		}
	}
	return true
}

func extractRegoSelectors(selectorsModule string) (selector, selector, error) {
	const regoMatchFuncMock = `
package library.v1.utils.labels.match.v1

all(labels, incl, excl) {
	true
}
`
	compiler, err := ast.CompileModules(map[string]string{
		"selectors.rego": selectorsModule,
		"mock.rego":      regoMatchFuncMock,
	})

	if err != nil {
		return nil, nil, err
	}

	include := extractSelector(compiler, "include")
	exclude := extractSelector(compiler, "exclude")

	return include, exclude, nil
}

func extractSelector(compiler *ast.Compiler, selectorName string) selector {
	rule := getRegoRuleByName(compiler, "selectors.rego", "systems")
	for _, expr := range rule.Body {
		values, ok := getSelectorValues(compiler, expr, selectorName)
		if ok {
			return values
		}
	}

	return nil
}

func getRegoRuleByName(compiler *ast.Compiler, moduleName, ruleName string) *ast.Rule {
	module, exists := compiler.Modules[moduleName]
	if !exists {
		return nil
	}

	for _, rule := range module.Rules {
		if rule.Head.Name.String() == ruleName {
			return rule
		}
	}

	return nil
}

func getSelectorValues(compiler *ast.Compiler, expr *ast.Expr, selectorName string) (selector, bool) {
	terms, ok := expr.Terms.([]*ast.Term)
	if !ok {
		return nil, false
	}

	if len(terms) != 3 || !isVarName(compiler, terms[1], selectorName) {
		return nil, false
	}

	obj, ok := terms[2].Value.(ast.Object)
	if !ok {
		return nil, false
	}

	selectorValues := make(map[string][]string)

	obj.Foreach(func(k, v *ast.Term) {
		key, ok := k.Value.(ast.String)
		if !ok {
			return
		}

		vSet, ok := v.Value.(ast.Set)
		if !ok {
			return
		}

		keyUnquote, err := strconv.Unquote(key.String())
		if err != nil {
			return
		}

		values := []string{}
		vSet.Foreach(func(e *ast.Term) {
			str, ok := e.Value.(ast.String)
			if !ok {
				return
			}
			name, err := strconv.Unquote(str.String())
			if err == nil {
				values = append(values, name)
			}
		})
		selectorValues[keyUnquote] = values
	})

	return selectorValues, true
}

func isVarName(compiler *ast.Compiler, term *ast.Term, varName string) bool {
	variable, ok := term.Value.(ast.Var)
	if !ok {
		return false
	}

	rewritten := getRewrittenVar(compiler, variable)
	return ok && rewritten.String() == varName
}

func getRewrittenVar(compiler *ast.Compiler, varName ast.Var) ast.Var {
	rewritten, exists := compiler.RewrittenVars[varName]
	if !exists {
		return varName
	}

	return rewritten
}
