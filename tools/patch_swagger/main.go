package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type formDataParam struct {
	Name        string
	Type        string
	Description string
	Required    bool
}

type patchRule struct {
	Path     string
	Method   string
	Consumes []string
	Params   []formDataParam
}

var patchRules = []patchRule{
	{
		Path:     "/v1/tenant/crm/candidates-resume",
		Method:   "post",
		Consumes: []string{"multipart/form-data"},
		Params: []formDataParam{
			{
				Name:        "file",
				Type:        "file",
				Description: "File to upload",
				Required:    true,
			},
		},
	},
}

func main() {
	inPath := flag.String("in", "", "path to input swagger.json")
	outPath := flag.String("out", "", "path to output swagger.json")
	flag.Parse()

	if *inPath == "" {
		exitf("flag -in is required")
	}
	if *outPath == "" {
		exitf("flag -out is required")
	}

	data, err := os.ReadFile(*inPath)
	if err != nil {
		exitf("read input file: %v", err)
	}

	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		exitf("unmarshal swagger json: %v", err)
	}

	if err := patchSwagger(doc, patchRules); err != nil {
		exitf("patch swagger: %v", err)
	}

	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		exitf("marshal patched swagger: %v", err)
	}
	out = append(out, '\n')

	if err := os.WriteFile(*outPath, out, 0o644); err != nil {
		exitf("write output file: %v", err)
	}
}

func patchSwagger(doc map[string]any, rules []patchRule) error {
	paths, ok := doc["paths"].(map[string]any)
	if !ok {
		return fmt.Errorf(`swagger does not contain object field "paths"`)
	}

	for _, rule := range rules {
		if err := applyPatchRule(paths, rule); err != nil {
			return err
		}
	}

	return nil
}

func applyPatchRule(paths map[string]any, rule patchRule) error {
	pathItemRaw, ok := paths[rule.Path]
	if !ok {
		return fmt.Errorf("path %q not found", rule.Path)
	}

	pathItem, ok := pathItemRaw.(map[string]any)
	if !ok {
		return fmt.Errorf("path item %q is not an object", rule.Path)
	}

	method := strings.ToLower(rule.Method)

	opRaw, ok := pathItem[method]
	if !ok {
		return fmt.Errorf("method %q for path %q not found", method, rule.Path)
	}

	op, ok := opRaw.(map[string]any)
	if !ok {
		return fmt.Errorf("operation %s %s is not an object", method, rule.Path)
	}

	op["consumes"] = toAnySlice(rule.Consumes)

	params := getParameters(op)

	filtered := make([]any, 0, len(params)+len(rule.Params))
	for _, p := range params {
		pm, ok := p.(map[string]any)
		if !ok {
			filtered = append(filtered, p)
			continue
		}

		inVal, _ := pm["in"].(string)
		nameVal, _ := pm["name"].(string)

		if inVal == "body" {
			continue
		}

		if inVal == "formData" && containsParam(rule.Params, nameVal) {
			continue
		}

		filtered = append(filtered, pm)
	}

	for _, p := range rule.Params {
		filtered = append(filtered, map[string]any{
			"name":        p.Name,
			"in":          "formData",
			"description": p.Description,
			"required":    p.Required,
			"type":        p.Type,
		})
	}

	op["parameters"] = filtered
	return nil
}

func getParameters(op map[string]any) []any {
	raw, ok := op["parameters"]
	if !ok {
		return nil
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil
	}
	return arr
}

func containsParam(params []formDataParam, name string) bool {
	for _, p := range params {
		if p.Name == name {
			return true
		}
	}
	return false
}

func toAnySlice(items []string) []any {
	out := make([]any, 0, len(items))
	for _, item := range items {
		out = append(out, item)
	}
	return out
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
