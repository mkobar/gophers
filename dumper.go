package gophers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

var sep = []byte("\r\n\r\n")

func isChunked(te []string) bool {
	for _, v := range te {
		if v == "chunked" {
			return true
		}
	}
	return false
}

func dump(b []byte, te []string) (headers []byte, body []byte, err error) {
	p := bytes.SplitN(b, sep, 2)
	headers, body = p[0], p[1]

	if len(body) > 0 && isChunked(te) {
		r := httputil.NewChunkedReader(bytes.NewReader(body))
		body, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, nil, err
		}
	}

	if len(body) > 0 {
		// TODO check that's really JSON
		var dst bytes.Buffer
		err = json.Indent(&dst, body, "", "  ")
		// TODO sort JSON object fields
		body = dst.Bytes()
	}

	return
}

func DumpRequest(req *http.Request) (headers []byte, body []byte, err error) {
	var b []byte
	b, err = httputil.DumpRequestOut(req, true)
	if err != nil {
		return
	}
	return dump(b, req.TransferEncoding)
}

func DumpResponse(res *http.Response) (headers []byte, body []byte, err error) {
	var b []byte
	b, err = httputil.DumpResponse(res, true)
	if err != nil {
		return
	}
	return dump(b, res.TransferEncoding)
}