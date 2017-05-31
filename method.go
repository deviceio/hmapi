package hmapi

var (
	DELETE  = &method{"DELETE"}
	GET     = &method{"GET"}
	HEAD    = &method{"HEAD"}
	OPTIONS = &method{"OPTIONS"}
	PATCH   = &method{"PATCH"}
	POST    = &method{"POST"}
	PUT     = &method{"PUT"}
)

type method struct {
	string
}

func (t method) String() string {
	return t.string
}
