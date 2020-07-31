package json

import (
	"encoding/json"

	"github.com/aintsashqa/link-shortener/src/shortener"
	"github.com/pkg/errors"
)

type RedirectJsonSerializer struct{}

func (r *RedirectJsonSerializer) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}

func (r *RedirectJsonSerializer) Encode(redirect *shortener.Redirect) ([]byte, error) {
	msg, err := json.Marshal(redirect)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return msg, nil
}
