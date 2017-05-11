package helpers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	//"github.com/pquerna/ffjson/ffjson"
	"html/template"
	"io"

	"github.com/ajg/form"
	"github.com/axiomzen/zenauth/helpers/mime"
)

type (
	// decoderFunc is a function that decodes a io.Reader into an interface
	decoderFunc func(v interface{}, r io.Reader) error
	// encoderFunc is a function that encodes a struct into the response writer
	encoderFunc func(v interface{}, w io.Writer) error
	// marshalFunc lets see how well we can wrap this
	marshalFunc func(v interface{}) ([]byte, error)
	// unmarshalFunc is the opposite of marshalFunc
	unmarshalFunc func(data []byte, v interface{}) error
)

var (
	// json decoding
	jsonDecoder = func(v interface{}, body io.Reader) error {
		//decoder := ffjson.NewDecoder()
		//return decoder.DecodeReader(body, v)

		decoder := json.NewDecoder(body)
		return decoder.Decode(v)
	}

	// xml decoding
	xmlDecoder = func(v interface{}, body io.Reader) error {
		decoder := xml.NewDecoder(body)
		return decoder.Decode(v)
	}

	// form decoding
	formDecoder = func(v interface{}, body io.Reader) error {
		decoder := form.NewDecoder(body)
		return decoder.Decode(v)
	}

	// supported decodings
	// add to this map if you want to support more types (like protobuffers, flatbuffers, thrift, etc)
	decodingMap = map[string]decoderFunc{
		"application/json":                     jsonDecoder,
		"application/xml":                      xmlDecoder,
		"application/vnd.google-earth.kml+xml": xmlDecoder,
		"application/x-www-form-urlencoded":    formDecoder,
	}

	unmarshalMap = map[string]unmarshalFunc{
		"application/json":                     json.Unmarshal,
		"application/xml":                      xml.Unmarshal,
		"application/vnd.google-earth.kml+xml": xml.Unmarshal,
		// we would have to wrap this one
		//"application/x-www-form-urlencoded":    form.Unmarshal,
	}

	jsonEncoder = func(v interface{}, w io.Writer) error {
		encoder := json.NewEncoder(w)
		return encoder.Encode(v)
	}

	// xmlMarshal swap out your xml marshaller here
	xmlMarshal = func(v interface{}) ([]byte, error) {
		return xml.MarshalIndent(v, "  ", "    ")
	}

	xmlEncoder = func(v interface{}, w io.Writer) error {
		encoder := xml.NewEncoder(w)
		return encoder.Encode(v)
	}

	formEncoder = func(v interface{}, w io.Writer) error {
		encoder := form.NewEncoder(w)
		return encoder.Encode(v)
	}

	tmpl = template.Must(template.New("text/html").Funcs(template.FuncMap{

		"RenderJSON": func(value interface{}) string {
			// render json pretty
			by, err := json.MarshalIndent(value, "", "    ")
			if err != nil {
				return err.Error()
			}
			// for now, I am skipping the syntax highlighting
			return string(by)
		},
	}).Delims("[[", "]]").Parse(
		`<html>
		<head>
		<meta http-equiv="content-type" content="text/html; charset=UTF-8">
		<meta name="robots" content="noindex, nofollow">
		<meta name="googlebot" content="noindex, nofollow">
		<style type="text/css">
		pre {
			background-color: ghostwhite;
			border: 1px solid silver;
			padding: 10px 20px;
			margin: 20px;
		}
		.json-key {
			color: teal;
		}
		.json-value {
			color: navy;
		}
		.json-string {
			color: brown;
		}
		</style>
		</head>
		<body>
		<pre>
		<code>
		[[ . | RenderJSON ]]
		</code>
		</pre>
		</body>
		</html>`))

	htmlEncoder = func(v interface{}, w io.Writer) error {

		str, ok := v.(string)
		if !ok {
			var buffer bytes.Buffer
			// render template
			if err := tmpl.Execute(&buffer, v); err != nil {
				return err
			}
			str = buffer.String()
		}
		_, err := w.Write([]byte(str))
		return err
	}

	htmlMarshal = func(v interface{}) ([]byte, error) {
		str, ok := v.(string)
		if ok {
			return []byte(str), nil
		}
		var buffer bytes.Buffer
		// render template
		if err := tmpl.Execute(&buffer, v); err != nil {
			return nil, err
		}
		return buffer.Bytes(), nil
	}

	plainEncoder = func(v interface{}, w io.Writer) error {
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%#v", v)
		}
		_, err := w.Write([]byte(str))
		return err
	}

	plainMarshal = func(v interface{}) ([]byte, error) {
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%#v", v)
		}

		return []byte(str), nil
	}

	encodingMap = map[string]encoderFunc{
		"":                                     jsonEncoder,
		"application/json":                     jsonEncoder,
		"application/xml":                      xmlEncoder,
		"application/vnd.google-earth.kml+xml": xmlEncoder,
		"application/x-www-form-urlencoded":    formEncoder,
		"text/html":                            htmlEncoder,
		"text/plain":                           plainEncoder,
	}

	marshalMap = map[string]marshalFunc{
		"":                                     json.Marshal,
		"application/json":                     json.Marshal,
		"application/xml":                      xmlMarshal,
		"application/vnd.google-earth.kml+xml": xmlMarshal,
		"text/html":                            htmlMarshal,
		"text/plain":                           htmlMarshal,
	}
)

// Decode is exposed so tests can re-use the decoding logic
// if we are decoding a body, then we must have been sent a content type
func Decode(v interface{}, contentType string, r io.Reader) error {
	//defer rc.Close()
	// apparently, The HTTP Client's Transport
	// is responsible for calling the Close method.

	// check quickly for a match
	if decFunc, ok := decodingMap[contentType]; ok {
		return decFunc(v, r)
	}

	// try cleaning it up
	ct, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	if decFunc, ok := decodingMap[ct]; ok {
		return decFunc(v, r)
	}

	return fmt.Errorf("unsupported content type: %s, %s", contentType, ct)
}

// Encode is exposed so tests can re-use the encoding logic
// encode must already know the content type as we wrote it
// in the header
func Encode(v interface{}, contentType string, w io.Writer) error {

	if encFunc, ok := encodingMap[contentType]; ok {
		return encFunc(v, w)
	}

	ct, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}

	if encFunc, ok := encodingMap[ct]; ok {
		return encFunc(v, w)
	}

	return fmt.Errorf("unsupported content type: %s, %s", contentType, ct)
}

// Marshal is exposed for testing, utility
// for marshall we must know a content type to dictate how to marshal
func Marshal(v interface{}, contentType string) ([]byte, error) {
	if mFunc, ok := marshalMap[contentType]; ok {
		return mFunc(v)
	}

	ct, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}
	if mFunc, ok := marshalMap[ct]; ok {
		return mFunc(v)
	}

	return nil, fmt.Errorf("unsupported content type: %s, %s", contentType, ct)
}

// Unmarshal is exposed for testing, utility
// will try to attempt to cleanup the contentType
func Unmarshal(data []byte, v interface{}, contentType string) error {
	if mFunc, ok := unmarshalMap[contentType]; ok {
		return mFunc(data, v)
	}

	ct, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}

	if mFunc, ok := unmarshalMap[ct]; ok {
		return mFunc(data, v)
	}

	return fmt.Errorf("unsupported content type: %s, %s", contentType, ct)
}
