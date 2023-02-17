package tls

import (
	_ "embed"
	"os"
)

//go:embed server.crt
var crt []byte

//go:embed server.key
var key []byte

func Refresh() error {
	f1, err := os.OpenFile("ca.crt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer f1.Close()
	_, err = f1.Write(crt)
	if err != nil {
		return err
	}

	f2, err := os.OpenFile("ca.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer f2.Close()
	_, err = f2.Write(key)
	if err != nil {
		return err
	}
	return nil
}
