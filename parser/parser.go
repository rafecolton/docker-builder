package parser

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/rafecolton/builder/builderfile"
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

func (me *Parser) Parse() (*builderfile.Builderfile, error) {
	bf := &builderfile.Builderfile{}

	raw, err := me.ParseRaw()
	if err != nil {
		return nil, err
	}

	if _, err := toml.Decode(raw, &bf); err != nil {
		return nil, err
	}

	return bf, nil
}

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
