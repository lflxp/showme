package decoder

import (
	"io"

	"github.com/juju/errors"
)

var DECODERS = map[string]Decoder{}

type Options struct {
	DeepDecode bool
}

type Decoder interface {
	Decode(io.Reader, io.Writer, *Options) error
	SetFilter(string)
}

func Register(name string, dec Decoder) {
	if _, ok := DECODERS[name]; !ok {
		DECODERS[name] = dec
	}
}

func GetDecoder(name string) (Decoder, error) {
	if dec, ok := DECODERS[name]; ok {
		return dec, nil
	}
	return nil, errors.New("Decoder not found: " + name)
}
