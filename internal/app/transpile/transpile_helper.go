package transpile

import (
	"bytes"
	"encoding/json"
	"log"

	ast "github.com/breml/logstash-config/ast"
)

type Constraints struct {
	Otherwise  bool
	Conditions []ast.Condition
}

func NewConstraint(_otherwise bool, _conditions ...ast.Condition) Constraints {
	return Constraints{
		Otherwise:  _otherwise,
		Conditions: _conditions,
	}
}

func NewConstraintLiteral() Constraints {
	return Constraints{
		Otherwise:  false,
		Conditions: []ast.Condition{},
	}
}

func AddCondToConstraint(c Constraints, otherwise bool, cond ast.Condition) Constraints {
	newC := Constraints{
		Otherwise:  otherwise,
		Conditions: append(c.Conditions, cond),
	}
	return newC
}

type IngestPipeline struct {
	Description         string            `json:"description"`
	Processors          []IngestProcessor `json:"processors"`
	OnFailureProcessors []IngestProcessor `json:"on_failure"`
}

func NewIngestPipeline() IngestPipeline {
	return IngestPipeline{
		Description:         "",
		Processors:          []IngestProcessor{},
		OnFailureProcessors: nil,
	}
}

// func processorsArrayToArrayMap(ips []IngestProcessor) []map[string]IngestProcessor {
// 	output := []map[string]IngestProcessor{}

// 	for _, ip := range ips {
// 		output = append(output, ToOutputMap(ip))
// 	}

// 	return output
// }

func MyJsonEncoder(m map[string]interface{}) string {
	buf := new(bytes.Buffer)

	e := json.NewEncoder(buf)

	e.SetEscapeHTML(false)

	err := e.Encode(m)

	if err != nil {
		// Panic
		log.Panicf("Something wrong in marshalling")
	}
	return string(buf.Bytes())
}

func (ingestPipeline IngestPipeline) String() string {

	m := map[string]interface{}{
		"description": ingestPipeline.Description,
		"processors":  ingestPipeline.Processors,
	}
	if ingestPipeline.OnFailureProcessors != nil {
		m["on_failure"] = ingestPipeline.OnFailureProcessors
	}
	return MyJsonEncoder(m)
}

type IngestProcessor interface {
	String() string
	IngestProcessorType() string
	SetIf(s *string) IngestProcessor
	// SetTag(string)
}

type SetProcessor struct {
	Value            string            `json:"value,omitempty"`
	Field            string            `json:"field"`
	CopyFrom         string            `json:"copy_from,omitempty"`
	Override         bool              `json:"override,omitempty"`
	IgnoreEmptyValue bool              `json:"ignore_empty_value,omitempty"`
	MediaType        string            `json:"media_type,omitempty"`
	Description      *string           `json:"description,omitempty"`
	If               *string           `json:"if,omitempty"`
	IgnoreFailure    bool              `json:"ignore_failure,omitempty"`
	Tag              string            `json:"tag"`
	OnFailure        []IngestProcessor `json:"on_failure,omitempty"`
}

func (ip SetProcessor) MarshalJSON() ([]byte, error) {
	type SetProcessorAlias SetProcessor

	return json.Marshal(
		map[string]SetProcessorAlias{
			ip.IngestProcessorType(): (SetProcessorAlias)(ip),
		},
	)
}

func StringHelper(ip IngestProcessor) string {

	return MyJsonEncoder(map[string]interface{}{
		ip.IngestProcessorType(): ip,
	})
}

func (sp SetProcessor) String() string {
	return StringHelper(sp)
}

func (sp SetProcessor) IngestProcessorType() string {
	return "set"
}

func (sp SetProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

type RemoveProcessor struct {
	Field         string            `json:"field"`
	IgnoreMissing bool              `json:"ignore_missing"`
	Keep          *string           `json:"keep"`
	Description   *string           `json:"description,omitempty"`
	If            *string           `json:"if,omitempty"`
	IgnoreFailure bool              `json:"ignore_failure,omitempty"`
	Tag           string            `json:"tag"`
	OnFailure     []IngestProcessor `json:"on_failure,omitempty"`
}

func (ip RemoveProcessor) String() string {
	return StringHelper(ip)
}

func (sp RemoveProcessor) IngestProcessorType() string {
	return "remove"
}

func (sp RemoveProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

type RenameProcessor struct {
	Field         string            `json:"field"`
	TargetField   string            `json:"target_field"`
	IgnoreMissing bool              `json:"ignore_missing"`
	Description   *string           `json:"description,omitempty"`
	If            *string           `json:"if,omitempty"`
	IgnoreFailure bool              `json:"ignore_failure,omitempty"`
	Tag           string            `json:"tag"`
	OnFailure     []IngestProcessor `json:"on_failure,omitempty"`
}

func (sp RenameProcessor) String() string {
	return StringHelper(sp)
}

func (sp RenameProcessor) IngestProcessorType() string {
	return "rename"
}

func (sp RenameProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

func (ip RenameProcessor) MarshalJSON() ([]byte, error) {
	type RenameProcessorAlias RenameProcessor

	return json.Marshal(
		map[string]RenameProcessorAlias{
			ip.IngestProcessorType(): (RenameProcessorAlias)(ip),
		},
	)
}

// Type for lowercase/uppercase
type CaseProcessor struct {
	Type          string            `json:"-"` // The field is only used internally to distinguish lowercase/uppercase
	Field         string            `json:"field"`
	TargetField   string            `json:"target_field"`
	IgnoreMissing bool              `json:"ignore_missing,omitempty"`
	Description   *string           `json:"description,omitempty"`
	If            *string           `json:"if,omitempty"`
	IgnoreFailure bool              `json:"ignore_failure,omitempty"`
	Tag           string            `json:"tag"`
	OnFailure     []IngestProcessor `json:"on_failure,omitempty"`
}

func (cp CaseProcessor) String() string {
	return StringHelper(cp)
}

func (cp CaseProcessor) IngestProcessorType() string {
	return cp.Type
}

func (sp CaseProcessor) ToOutputMap() map[string]interface{} {
	return map[string]interface{}{}
}

func (sp CaseProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

func (ip CaseProcessor) MarshalJSON() ([]byte, error) {
	type CaseProcessorAlias CaseProcessor

	return json.Marshal(
		map[string]CaseProcessorAlias{
			ip.IngestProcessorType(): (CaseProcessorAlias)(ip),
		},
	)
}

// Type for lowercase/uppercase
type GrokProcessor struct {
	Field              string            `json:"field,omitempty"`
	Patterns           []string          `json:"patterns,omitempty"`
	PatternDefinitions map[string]string `json:"pattern_definitions,omitempty"`
	ECSCompatibility   string            `json:"ecs_compatibility,omitempty"`
	TraceMatch         bool              `json:"trace_match,omitempty"`
	IgnoreMissing      bool              `json:"ignore_missing,omitempty"`
	Description        *string           `json:"description,omitempty"`
	If                 *string           `json:"if,omitempty"`
	IgnoreFailure      bool              `json:"ignore_failure,omitempty"`
	Tag                string            `json:"tag"`
	OnFailure          []IngestProcessor `json:"on_failure,omitempty"`
}

func (gp GrokProcessor) String() string {
	return StringHelper(gp)
}

func (gp GrokProcessor) IngestProcessorType() string {
	return "grok"
}

func (sp GrokProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

func (ip GrokProcessor) MarshalJSON() ([]byte, error) {
	type GrokProcessorAlias GrokProcessor

	return json.Marshal(
		map[string]GrokProcessorAlias{
			ip.IngestProcessorType(): (GrokProcessorAlias)(ip),
		},
	)
}

type AppendProcessor struct {
	Field           string            `json:"field,omitempty"`
	Value           []string          `json:"value,omitempty"`
	AllowDuplicates bool              `json:"allow_duplicates,omitempty"`
	MediaType       *string           `json:"media_type,omitempty"`
	Description     *string           `json:"description,omitempty"`
	If              *string           `json:"if,omitempty"`
	IgnoreFailure   bool              `json:"ignore_failure,omitempty"`
	Tag             string            `json:"tag"`
	OnFailure       []IngestProcessor `json:"on_failure,omitempty"`
}

func (ap AppendProcessor) String() string {
	return StringHelper(ap)
}

func (ap AppendProcessor) IngestProcessorType() string {
	return "append"
}

func (ip AppendProcessor) MarshalJSON() ([]byte, error) {
	type AppendProcessorAlias AppendProcessor

	return json.Marshal(
		map[string]AppendProcessorAlias{
			ip.IngestProcessorType(): (AppendProcessorAlias)(ip),
		},
	)
}

func (sp AppendProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

type GsubProcessor struct {
	Field         string            `json:"field,omitempty"`
	Pattern       string            `json:"patterns,omitempty"`
	Replacement   string            `json:"replacement"`
	TargetField   *string           `json:"target_field,omitempty"`
	IgnoreMissing bool              `json:"ignore_missing,omitempty"`
	Description   *string           `json:"description,omitempty"`
	If            *string           `json:"if,omitempty"`
	IgnoreFailure bool              `json:"ignore_failure,omitempty"`
	Tag           string            `json:"tag"`
	OnFailure     []IngestProcessor `json:"on_failure,omitempty"`
}

func (ap GsubProcessor) String() string {
	return StringHelper(ap)
}

func (ap GsubProcessor) IngestProcessorType() string {
	return "gsub"
}

func (sp GsubProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

func (ip GsubProcessor) MarshalJSON() ([]byte, error) {
	type GsubProcessorAlias GsubProcessor

	return json.Marshal(
		map[string]GsubProcessorAlias{
			ip.IngestProcessorType(): (GsubProcessorAlias)(ip),
		},
	)
}

type JoinProcessor struct {
	Field         string            `json:"field,omitempty"`
	Separator     string            `json:"separator,omitempty"`
	TargetField   *string           `json:"target_field,omitempty"`
	IgnoreMissing bool              `json:"ignore_missing,omitempty"`
	Description   *string           `json:"description,omitempty"`
	If            *string           `json:"if,omitempty"`
	IgnoreFailure bool              `json:"ignore_failure,omitempty"`
	Tag           string            `json:"tag"`
	OnFailure     []IngestProcessor `json:"on_failure,omitempty"`
}

func (ap JoinProcessor) String() string {
	return StringHelper(ap)
}

func (ap JoinProcessor) IngestProcessorType() string {
	return "join"
}

func (sp JoinProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

func (ip JoinProcessor) MarshalJSON() ([]byte, error) {
	type JoinProcessorAlias JoinProcessor

	return json.Marshal(
		map[string]JoinProcessorAlias{
			ip.IngestProcessorType(): (JoinProcessorAlias)(ip),
		},
	)
}

type KVProcessor struct {
	Field         string            `json:"field,omitempty"`
	FieldSplit    string            `json:"field_split,omitempty"`
	ValueSplit    string            `json:"value_split,omitempty"`
	TargetField   *string           `json:"target_field,omitempty"`
	IncludeKeys   []string          `json:"include_keys,omitempty"`
	ExcludeKeys   []string          `json:"exclude_keys,omitempty"`
	IgnoreMissing bool              `json:"ignore_missing,omitempty"`
	Prefix        *string           `json:"prefix,omitempty"`
	TrimKey       *string           `json:"trim_key,omitempty"`
	TrimValue     *string           `json:"trim_value,omitempty"`
	StripBrackets bool              `json:"strip_bracket,omitempty"`
	Pattern       string            `json:"patterns,omitempty"`
	Description   *string           `json:"description,omitempty"`
	If            *string           `json:"if,omitempty"`
	IgnoreFailure bool              `json:"ignore_failure,omitempty"`
	Tag           string            `json:"tag"`
	OnFailure     []IngestProcessor `json:"on_failure,omitempty"`
}

func (ap KVProcessor) String() string {
	return StringHelper(ap)
}

func (ap KVProcessor) IngestProcessorType() string {
	return "kv"
}

func (sp KVProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

func (ip KVProcessor) MarshalJSON() ([]byte, error) {
	type KVProcessorAlias KVProcessor

	return json.Marshal(
		map[string]KVProcessorAlias{
			ip.IngestProcessorType(): (KVProcessorAlias)(ip),
		},
	)
}

type DissectProcessor struct {
	Field           string            `json:"field,omitempty"`
	Pattern         string            `json:"patterns,omitempty"`
	AppendSeparator *string           `json:"append_separator,omitempty"`
	IgnoreMissing   bool              `json:"ignore_missing,omitempty"`
	Description     *string           `json:"description,omitempty"`
	If              *string           `json:"if,omitempty"`
	IgnoreFailure   bool              `json:"ignore_failure,omitempty"`
	Tag             string            `json:"tag"`
	OnFailure       []IngestProcessor `json:"on_failure,omitempty"`
}

func (ap DissectProcessor) String() string {
	return StringHelper(ap)
}

func (ap DissectProcessor) IngestProcessorType() string {
	return "dissect"
}

func (sp DissectProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

func (ip DissectProcessor) MarshalJSON() ([]byte, error) {
	type DissectProcessorAlias DissectProcessor

	return json.Marshal(
		map[string]DissectProcessorAlias{
			ip.IngestProcessorType(): (DissectProcessorAlias)(ip),
		},
	)
}

type DateProcessor struct {
	Field         string            `json:"field,omitempty"`
	TargetField   *string           `json:"target_field,omitempty"` // Default value @timestamp
	Formats       []string          `json:"formats,omitempty"`
	Pattern       string            `json:"patterns,omitempty"`
	Timezone      *string           `json:"timezone,omitempty"`
	Locale        *string           `json:"locale,omitempty"`
	OutputFormat  *string           `json:"output_format,omitempty"` // Default yyyy-MM-dd'T'HH:mm:ss.SSSXXX
	IgnoreMissing bool              `json:"ignore_missing,omitempty"`
	Description   *string           `json:"description,omitempty"`
	If            *string           `json:"if,omitempty"`
	IgnoreFailure bool              `json:"ignore_failure,omitempty"`
	Tag           string            `json:"tag"`
	OnFailure     []IngestProcessor `json:"on_failure,omitempty"`
}

func (ap DateProcessor) String() string {
	return StringHelper(ap)
}

func (ap DateProcessor) IngestProcessorType() string {
	return "dissect"
}

func (sp DateProcessor) SetIf(s *string) IngestProcessor {
	sp.If = s
	return sp
}

func (ip DateProcessor) MarshalJSON() ([]byte, error) {
	type DateProcessorAlias DateProcessor

	return json.Marshal(
		map[string]DateProcessorAlias{
			ip.IngestProcessorType(): (DateProcessorAlias)(ip),
		},
	)
}
