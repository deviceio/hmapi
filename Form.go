package hmapi

type Form struct {
	Action  string       `json:"action,omitempty"`
	Method  string       `json:"method"`
	Type    MediaType    `json:"type,omitempty"`
	Enctype MediaType    `json:"enctype,omitempty"`
	Fields  []*FormField `json:"fields,omitempty"`
}
