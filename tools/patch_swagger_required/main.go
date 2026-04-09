package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	validate "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type fieldInfo struct {
	ProtoName string
	JSONName  string

	Optional bool
	Required bool
}

type messageInfo struct {
	FullName  string
	ShortName string

	FieldsByProtoName map[string]fieldInfo
	FieldsByJSONName  map[string]fieldInfo
}

func main() {
	inPath := flag.String("in", "", "path to input swagger.json")
	descriptorPath := flag.String("descriptor", "", "path to FileDescriptorSet (proto.pb)")
	outPath := flag.String("out", "", "path to output swagger.json")
	flag.Parse()

	if *inPath == "" {
		exitf("flag -in is required")
	}
	if *descriptorPath == "" {
		exitf("flag -descriptor is required")
	}
	if *outPath == "" {
		exitf("flag -out is required")
	}

	doc, err := readSwagger(*inPath)
	if err != nil {
		exitf("read swagger: %v", err)
	}

	fds, err := loadDescriptorSet(*descriptorPath)
	if err != nil {
		exitf("load descriptor set: %v", err)
	}

	messages := collectMessages(fds)

	if err := patchSwaggerRequired(doc, messages); err != nil {
		exitf("patch swagger required: %v", err)
	}

	if err := writeSwagger(*outPath, doc); err != nil {
		exitf("write swagger: %v", err)
	}
}

func readSwagger(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func writeSwagger(path string, doc map[string]any) error {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	return os.WriteFile(path, data, 0o644)
}

func loadDescriptorSet(path string) (*descriptorpb.FileDescriptorSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var fds descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &fds); err != nil {
		return nil, err
	}

	return &fds, nil
}

func patchSwaggerRequired(doc map[string]any, messages []messageInfo) error {
	definitionsRaw, ok := doc["definitions"]
	if !ok {
		return nil
	}

	definitions, ok := definitionsRaw.(map[string]any)
	if !ok {
		return fmt.Errorf(`field "definitions" is not an object`)
	}

	for _, defRaw := range definitions {
		def, ok := defRaw.(map[string]any)
		if !ok {
			continue
		}

		propertiesRaw, ok := def["properties"]
		if !ok {
			continue
		}

		properties, ok := propertiesRaw.(map[string]any)
		if !ok || len(properties) == 0 {
			continue
		}

		msg, found := findBestMessageForDefinition(properties, messages)
		if !found {
			continue
		}

		required := make([]string, 0, len(properties))

		for swaggerFieldName := range properties {
			fi, ok := lookupField(msg, swaggerFieldName)
			if !ok {
				continue
			}

			if fi.Required {
				required = append(required, swaggerFieldName)
				continue
			}

			if !fi.Optional {
				required = append(required, swaggerFieldName)
			}
		}

		sort.Strings(required)

		if len(required) == 0 {
			delete(def, "required")
			continue
		}

		def["required"] = toAnySlice(required)
	}

	return nil
}

func findBestMessageForDefinition(
	properties map[string]any,
	messages []messageInfo,
) (messageInfo, bool) {
	var matched messageInfo
	bestScore := -1

	for _, msg := range messages {
		score := scoreMessageMatch(msg, properties)
		if score > bestScore {
			bestScore = score
			matched = msg
		}
	}

	if bestScore <= 0 {
		return messageInfo{}, false
	}

	return matched, true
}

func scoreMessageMatch(msg messageInfo, properties map[string]any) int {
	score := 0
	matchedFields := 0

	for swaggerFieldName := range properties {
		if _, ok := msg.FieldsByJSONName[swaggerFieldName]; ok {
			score += 10
			matchedFields++
			continue
		}
		if _, ok := msg.FieldsByProtoName[swaggerFieldName]; ok {
			score += 7
			matchedFields++
			continue
		}
	}

	// штрафуем большие message с кучей лишних полей,
	// чтобы более точный message побеждал при одинаковом числе совпадений.
	totalFields := len(msg.FieldsByJSONName)
	extraFields := totalFields - matchedFields
	if extraFields > 0 {
		score -= extraFields
	}

	// бонус за полное совпадение по количеству полей
	if matchedFields == len(properties) && totalFields == len(properties) {
		score += 1000
	}

	return score
}

func lookupField(msg messageInfo, swaggerFieldName string) (fieldInfo, bool) {
	if fi, ok := msg.FieldsByJSONName[swaggerFieldName]; ok {
		return fi, true
	}
	if fi, ok := msg.FieldsByProtoName[swaggerFieldName]; ok {
		return fi, true
	}
	return fieldInfo{}, false
}

func collectMessages(fds *descriptorpb.FileDescriptorSet) []messageInfo {
	out := make([]messageInfo, 0)

	for _, file := range fds.File {
		pkg := file.GetPackage()
		for _, msg := range file.GetMessageType() {
			walkMessage(pkg, nil, msg, &out)
		}
	}

	return out
}

func walkMessage(
	pkg string,
	parent []string,
	msg *descriptorpb.DescriptorProto,
	out *[]messageInfo,
) {
	fullParts := make([]string, 0, len(parent)+1)
	fullParts = append(fullParts, parent...)
	fullParts = append(fullParts, msg.GetName())

	shortName := msg.GetName()
	fullName := shortName
	if len(fullParts) > 1 {
		fullName = strings.Join(fullParts, ".")
	}
	if pkg != "" {
		fullName = pkg + "." + fullName
	}

	mi := messageInfo{
		FullName:          fullName,
		ShortName:         shortName,
		FieldsByProtoName: make(map[string]fieldInfo),
		FieldsByJSONName:  make(map[string]fieldInfo),
	}

	oneofIndexes := map[int32]struct{}{}
	for i := range msg.GetOneofDecl() {
		oneofIndexes[int32(i)] = struct{}{}
	}

	for _, f := range msg.GetField() {
		hasDecision, required := getExplicitRequiredDecision(f)

		optional := false
		if hasDecision {
			optional = !required
		} else {
			optional = isFieldOptional(f, oneofIndexes)
		}

		fi := fieldInfo{
			ProtoName: f.GetName(),
			JSONName:  effectiveJSONName(f),
			Optional:  optional,
			Required:  required,
		}

		mi.FieldsByProtoName[fi.ProtoName] = fi
		mi.FieldsByJSONName[fi.JSONName] = fi
	}

	*out = append(*out, mi)

	nextParent := append(append([]string(nil), parent...), msg.GetName())
	for _, nested := range msg.GetNestedType() {
		if nested.GetOptions().GetMapEntry() {
			continue
		}
		walkMessage(pkg, nextParent, nested, out)
	}
}

func effectiveJSONName(f *descriptorpb.FieldDescriptorProto) string {
	if v := f.GetJsonName(); v != "" {
		return v
	}
	return protoToJSONName(f.GetName())
}

func protoToJSONName(s string) string {
	if s == "" {
		return s
	}

	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}

	var b strings.Builder
	b.WriteString(parts[0])

	for _, p := range parts[1:] {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}

	return b.String()
}

func isFieldOptional(
	f *descriptorpb.FieldDescriptorProto,
	oneofIndexes map[int32]struct{},
) bool {
	if f.GetProto3Optional() {
		return true
	}

	if _, ok := oneofIndexes[f.GetOneofIndex()]; ok && f.OneofIndex != nil {
		return true
	}

	if f.GetLabel() == descriptorpb.FieldDescriptorProto_LABEL_REPEATED {
		return true
	}

	if f.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		return true
	}

	return false
}

// getExplicitRequiredDecision returns:
//   - hasDecision=true, required=true  => field is explicitly required
//   - hasDecision=true, required=false => field is explicitly optional
//   - hasDecision=false                => fallback to heuristic
func getExplicitRequiredDecision(f *descriptorpb.FieldDescriptorProto) (bool, bool) {
	if has, required := getBufValidateRequiredDecision(f); has {
		return true, required
	}

	if isRequiredByGoogleFieldBehavior(f) {
		return true, true
	}

	return false, false
}

// explicit decision exists only when `required` itself is present.
func getBufValidateRequiredDecision(f *descriptorpb.FieldDescriptorProto) (bool, bool) {
	opts := f.GetOptions()
	if opts == nil {
		return false, false
	}

	if !proto.HasExtension(opts, validate.E_Field) {
		return false, false
	}

	ext := proto.GetExtension(opts, validate.E_Field)
	rules, ok := ext.(*validate.FieldRules)
	if !ok || rules == nil {
		return false, false
	}

	fd := rules.ProtoReflect().Descriptor().Fields().ByName("required")
	if fd == nil {
		return false, false
	}

	if !rules.ProtoReflect().Has(fd) {
		return false, false
	}

	return true, rules.ProtoReflect().Get(fd).Bool()
}

func isRequiredByGoogleFieldBehavior(f *descriptorpb.FieldDescriptorProto) bool {
	opts := f.GetOptions()
	if opts == nil {
		return false
	}

	if !proto.HasExtension(opts, annotations.E_FieldBehavior) {
		return false
	}

	ext := proto.GetExtension(opts, annotations.E_FieldBehavior)
	behaviors, ok := ext.([]annotations.FieldBehavior)
	if !ok {
		return false
	}

	for _, b := range behaviors {
		if b == annotations.FieldBehavior_REQUIRED {
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
