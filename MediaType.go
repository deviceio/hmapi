package hmapi

type MediaType string

func (t *MediaType) String() string {
	return string(*t)
}

var (
	MediaTypeHMAPIBoolean      = MediaType("application/vnd.hmapi.Bool;charset=utf-8")
	MediaTypeHMAPIFloat32      = MediaType("application/vnd.hmapi.Float32;charset=utf-8")
	MediaTypeHMAPIFloat64      = MediaType("application/vnd.hmapi.Float64;charset=utf-8")
	MediaTypeHMAPIInt          = MediaType("application/vnd.hmapi.Int;charset=utf-8")
	MediaTypeHMAPIInt32        = MediaType("application/vnd.hmapi.Int32;charset=utf-8")
	MediaTypeHMAPIInt64        = MediaType("application/vnd.hmapi.Int64;charset=utf-8")
	MediaTypeHMAPIString       = MediaType("application/vnd.hmapi.String;charset=utf-8")
	MediaTypeHMAPIUInt         = MediaType("application/vnd.hmapi.UInt;charset=utf-8")
	MediaTypeHMAPIUInt32       = MediaType("application/vnd.hmapi.UInt32;charset=utf-8")
	MediaTypeHMAPIUInt64       = MediaType("application/vnd.hmapi.UInt64;charset=utf-8")
	MediaTypeOctetStream       = MediaType("application/octet-stream")
	MediaTypeJSON              = MediaType("application/json;charset=utf-8")
	MediaTypeTextPlain         = MediaType("text/plain;charset=utf-8")
	MediaTypeMultipartFormData = MediaType("multipart/form-data")
)
