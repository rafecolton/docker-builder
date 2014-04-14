package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Parser struct {
	Builderfile string
}

func New() *Parser {
	p := &Parser{}

	return p
}

var (
	err error
)

func (me *Parser) ParseRaw() (string, error) {

	if !me.IsOpenable() {
		return "", errors.New(fmt.Sprintf("%s is not openable", me.Builderfile))
	}

	bytes, err := ioutil.ReadFile(me.Builderfile)

	if err != nil {
		return "", err
	}

	raw := string(bytes)

	return raw, nil
}

func (me *Parser) IsOpenable() bool {

	file, err := os.Open(me.Builderfile)
	defer file.Close()

	if err != nil {
		return false
	}

	return true
}
