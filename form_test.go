package hmapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Test_FormRequest_when_calling_submit struct {
	suite.Suite
}

func (t *Test_FormRequest_when_calling_submit) Test_multipart_form_successfully_submitted() {
	ret := t.getTestServerAndClient()
	defer ret.Server.Close()

	ret.Mux.HandleFunc("/resource/test", func(rw http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(4096)

		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			assert.Fail(t.T(), "error parsing form", err.Error())
			return
		}

		if r.Form.Get("foo") != "test" {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			assert.Fail(t.T(), "form field foo not supplied", "")
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("response"))
	}).Methods("POST")

	ret.Mux.HandleFunc("/resource", func(rw http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(&Resource{
			Forms: map[string]*Form{
				"test": &Form{
					Action:  "/resource/test",
					Method:  POST,
					Enctype: MediaTypeMultipartFormData,
					Type:    "none",
					Fields: []*FormField{
						&FormField{
							Name:     "foo",
							Type:     MediaTypeHMAPIString,
							Required: true,
						},
					},
				},
			},
		})

		if err != nil {
			assert.FailNow(t.T(), "error marshal json", err)
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(b)
	}).Methods("GET")

	submission := ret.Client.Resource("/resource").Form("test").AddFieldAsString("foo", "test").Submit()

	<-submission.Done()

	assert.Nil(t.T(), submission.Err())
	assert.NotNil(t.T(), submission.Response())
	assert.Equal(t.T(), http.StatusOK, submission.Response().StatusCode)

	respbody, _ := ioutil.ReadAll(submission.Response().Body)

	log.Println(respbody)

	assert.Equal(t.T(), "response", string(respbody))
}

func (t *Test_FormRequest_when_calling_submit) getTestServerAndClient() (ret struct {
	Mux    *mux.Router
	Host   string
	Port   int
	Server *httptest.Server
	Client Client
}) {
	mux := mux.NewRouter()
	http.Handle("/", mux)
	svr := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Println("TEST SVR REQUEST", r)
		mux.ServeHTTP(rw, r)
	}))

	url, _ := url.Parse(svr.URL)

	hoststr, portstr, _ := net.SplitHostPort(url.Host)
	port, _ := strconv.ParseInt(portstr, 10, 0)

	client := NewClient(&ClientConfig{
		Auth:   &AuthNone{},
		Host:   hoststr,
		Port:   int(port),
		Scheme: HTTP,
	})

	ret.Mux = mux
	ret.Host = hoststr
	ret.Port = int(port)
	ret.Server = svr
	ret.Client = client
	return
}

func TestRunFormTestSuites(t *testing.T) {
	suite.Run(t, new(Test_FormRequest_when_calling_submit))
}
