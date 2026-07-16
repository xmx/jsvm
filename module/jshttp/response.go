package jshttp

import (
	"encoding/json"
	"io"
	"net/http"
)

type response struct {
	*http.Response
}

func (r *response) JSON() (any, error) {
	var ret any
	dec := json.NewDecoder(r.Response.Body)
	if err := dec.Decode(&ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *response) Text() (string, error) {
	bs, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}
