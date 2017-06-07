package hmapi

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type FormRequest interface {
	AddField(name string, media MediaType, value interface{}) FormRequest
	AddFieldAsString(name string, value string) FormRequest
	AddFieldAsBool(name string, value bool) FormRequest
	AddFieldAsOctetStream(name string, value io.Reader) FormRequest
	AddFieldAsInt(name string, value int) FormRequest
	Submit() FormSubmission
}

type FormSubmission interface {
	Response() *http.Response
	Err() error
	Done() <-chan struct{}
	Cancel()
}

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

type formRequest struct {
	name     string
	fields   []*formField
	resource *resourceRequest
}

type formSubmission struct {
	httpResponse *http.Response
	err          error
	ctx          context.Context
	cancel       context.CancelFunc
}

func (t *formSubmission) Response() *http.Response {
	return t.httpResponse
}

func (t *formSubmission) Err() error {
	return t.err
}

func (t *formSubmission) Cancel() {
	t.cancel()
}

func (t *formSubmission) Done() <-chan struct{} {
	return t.ctx.Done()
}

type formField struct {
	name      string
	mediaType MediaType
	value     interface{}
}

func (t *formRequest) AddField(name string, media MediaType, value interface{}) FormRequest {
	t.fields = append(t.fields, &formField{
		name:      name,
		mediaType: media,
		value:     value,
	})

	return t
}

func (t *formRequest) AddFieldAsString(name string, value string) FormRequest {
	t.AddField(name, MediaTypeHMAPIString, value)
	return t
}

func (t *formRequest) AddFieldAsBool(name string, value bool) FormRequest {
	t.AddField(name, MediaTypeHMAPIBoolean, value)
	return t
}

func (t *formRequest) AddFieldAsOctetStream(name string, value io.Reader) FormRequest {
	t.AddField(name, MediaTypeOctetStream, value)
	return t
}

func (t *formRequest) AddFieldAsInt(name string, value int) FormRequest {
	t.AddField(name, MediaTypeHMAPIInt, value)
	return t
}

func (t *formRequest) Submit() FormSubmission {
	ctx, ctxcancel := context.WithCancel(context.Background())

	submission := &formSubmission{
		ctx:    ctx,
		cancel: ctxcancel,
	}

	go t.submit(submission)

	return submission
}

func (t *formRequest) submit(submission *formSubmission) {
	hmres, err := t.resource.Get()

	if err != nil {
		submission.err = err
		submission.cancel()
		return
	}

	hmform, ok := hmres.Forms[t.name]

	if !ok {
		submission.err = &ErrResourceNoSuchForm{
			FormName: t.name,
			Resource: t.resource.path,
		}
		submission.cancel()
		return
	}

	bodyr, bodyw := io.Pipe()

	request, err := http.NewRequest(
		hmform.Method.String(),
		t.resource.client.baseuri+hmform.Action,
		bodyr,
	)

	if err != nil {
		submission.err = err
		submission.cancel()
		return
	}

	request = request.WithContext(submission.ctx)

	switch hmform.Enctype {
	case MediaTypeMultipartFormData:
		request.Header.Set("Content-Type", MediaTypeMultipartFormData.String())
	default:
		submission.err = &ErrUnsupportedMediaType{
			MediaType: hmform.Enctype,
		}
		submission.cancel()
		return
	}

	chresp := make(chan *http.Response)
	chresperr := make(chan error)
	chformerr := make(chan error)

	go func() {
		resp, err := t.resource.client.do(request)
		chresperr <- err
		chresp <- resp
	}()

	go func() {
		switch hmform.Enctype {
		case MediaTypeMultipartFormData:
			t.writeMultipartForm(bodyw, hmform, submission, chformerr)
			bodyw.Close()
		}
	}()

waitforcomplete:
	for {
		select {
		case formerr := <-chformerr:
			submission.err = formerr
		case resperr := <-chresperr:
			submission.err = resperr
		case resp := <-chresp:
			submission.httpResponse = resp
			break waitforcomplete
		case <-submission.ctx.Done():
			break waitforcomplete
		}
	}

	submission.cancel() //done
}

func (t *formRequest) writeMultipartForm(writer io.Writer, form *Form, submission *formSubmission, cherr chan error) {
	mpwriter := multipart.NewWriter(writer)
	mpwriter.SetBoundary(MultipartFormDataBoundry)
	defer mpwriter.Close()

	for _, field := range t.fields {
		select {
		case <-submission.ctx.Done():
			return
		default:
		}

		switch field.mediaType {
		case MediaTypeOctetStream:
			fieldwriter, _ := mpwriter.CreateFormField(field.name)
			if _, err := io.Copy(fieldwriter, field.value.(io.Reader)); err != nil {
				cherr <- err
				return
			}

		case MediaTypeHMAPIInt:
			if err := mpwriter.WriteField(field.name, strconv.FormatInt(int64(field.value.(int)), 10)); err != nil {
				cherr <- err
				return
			}

		case MediaTypeHMAPIString:
			if err := mpwriter.WriteField(field.name, field.value.(string)); err != nil {
				cherr <- err
				return
			}

		case MediaTypeHMAPIBoolean:
			if err := mpwriter.WriteField(field.name, strconv.FormatBool(field.value.(bool))); err != nil {
				cherr <- err
				return
			}

		default:
			cherr <- &ErrUnsupportedMediaType{
				MediaType: form.Enctype,
			}
			return
		}
	}

	cherr <- nil
}

type formResponse struct {
	httpResponse *http.Response
}

func (t *formResponse) HttpResponse() *http.Response {
	return t.httpResponse
}
