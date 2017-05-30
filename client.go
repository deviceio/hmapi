package hmapi

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

type Client interface {
	Resource(path string) ResourceRequest
}

type client struct {
	baseuri    string
	httpclient *http.Client
	host       string
	port       int
	auth       Auth
}

func NewClient(scheme scheme, host string, port int, auth Auth) Client {
	if auth == nil {
		auth = &AuthNone{}
	}

	return &client{
		baseuri: fmt.Sprintf("%v://%v:%v", string(scheme), host, port),
		httpclient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		auth: auth,
	}
}

func (t *client) Resource(path string) ResourceRequest {
	return &resourceRequest{
		client: t,
		path:   path,
	}
}

func (t *client) do(r *http.Request) (*http.Response, error) {
	t.auth.Sign(r)
	return t.httpclient.Do(r)
}