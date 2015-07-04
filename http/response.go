package http

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"html/template"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

func NewResponse(rw http.ResponseWriter) *Response {
	return &Response{
		ResponseWriter: rw,
	}
}

func (r *Response) WriteJSON(v interface{}) (content []byte, err error) {
	content, err = json.Marshal(v)
	if err != nil {
		return nil, err
	}
	r.Header().Set("Content-Type", "application/json;charset=utf-8")
	_, err = r.Write(content)
	return
}

func (r *Response) WriteXML(v interface{}) (content []byte, err error) {
	content, err = xml.Marshal(v)
	if err != nil {
		return nil, err
	}
	r.Header().Set("Content-Type", "application/xml;charset=utf-8")
	_, err = r.Write([]byte(xml.Header))
	_, err = r.Write(content)
	return
}

func (r *Response) WriteJSONP(v interface{}, callback string) (content []byte, err error) {
	if callback == "" {
		err = errors.New("jsonp callback not found")
		return nil, err
	}
	r.Header().Set("Content-Type", "application/javascript;charset=utf-8")
	content, err = json.Marshal(v)
	if err != nil {
		return nil, err
	}

	content = append([]byte{'('}, content...)
	content = append(content, []byte{')'}...)
	content = append([]byte(template.JSEscapeString(callback)), content...)

	_, err = r.Write(content)
	return
}
