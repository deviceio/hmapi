package hmapi

type FormField struct {
	Name     string      `json:"name"`
	Type     MediaType   `json:"type,omitempty"`
	Encoding MediaType   `json:"encoding,omitempty"`
	Required bool        `json:"required"`
	Value    interface{} `json:"value,omitempty"`
}
