package parser

import (
	"github.com/rafecolton/docker-builder/builderfile"
	"github.com/rafecolton/docker-builder/parser/uuid"
)

import (
	"path/filepath"

	"github.com/Sirupsen/logrus"
)

/*
Parser is a struct that contains a Builderfile and knows how to parse it both
as raw text and to convert toml to a Builderfile struct.  It also knows how to
tell if the Builderfile is valid (openable) or nat.
*/
type Parser struct {
	filename string
	*logrus.Logger
	uuidGenerator uuid.Generator
	top           string
}

/*
NewParser returns an initialized Parser.  Not currently necessary, as no
default values are assigned to a new Parser, but useful to have in case we need
to change this.
*/
func NewParser(filename string, l *logrus.Logger) *Parser {
	builderfile.Logger(l)
	return &Parser{
		Logger:        l,
		filename:      filename,
		uuidGenerator: uuid.NewUUIDGenerator(true),
		top:           filepath.Dir(filename),
	}
}
