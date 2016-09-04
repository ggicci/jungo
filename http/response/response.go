package response

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"errors"
	"html/template"
	"net/http"
)

func WriteJSON(rw http.ResponseWriter, v interface{}) error {
	encoder := json.NewEncoder(rw)
	rw.Header().Set("Content-Type", "application/json")
	return encoder.Encode(v)
}

func WriteJSONP(rw http.ResponseWriter, v interface{}, callback string) error {
	if callback == "" {
		return errors.New("jsonp callback not found")
	}
	output := bufio.NewWriter(rw)
	encoder := json.NewEncoder(output)
	rw.Header().Set("Content-Type", "application/javascript;charset=utf-8")
	output.WriteString(template.JSEscapeString(callback))
	output.WriteByte('(')
	if err := encoder.Encode(v); err != nil {
		return err
	}
	output.WriteByte(')')
	return output.Flush()
}

func WriteXML(rw http.ResponseWriter, v interface{}) error {
	encoder := xml.NewEncoder(rw)
	rw.Header().Set("Content-Type", "application/xml;charset=utf-8")
	return encoder.Encode(v)
}
