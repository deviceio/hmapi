package hmapi

type Form struct {
	Action  string       `json:"action,omitempty"`
	Method  method       `json:"method"`
	Type    MediaType    `json:"type,omitempty"`
	Enctype MediaType    `json:"enctype,omitempty"`
	Fields  []*FormField `json:"fields,omitempty"`
}

type FormField struct {
	Name     string      `json:"name"`
	Type     MediaType   `json:"type,omitempty"`
	Encoding MediaType   `json:"encoding,omitempty"`
	Required bool        `json:"required"`
	Multiple bool        `json:"multiple"`
	Value    interface{} `json:"value,omitempty"`
}
