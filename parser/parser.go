package parser

import (
	"fmt"
	"github.com/rafecolton/bob/log"
	"github.com/rafecolton/bob/parser/dclient"
)

/*
Parser is a struct that contains a Builderfile and knows how to parse it both
as raw text and to convert toml to a Builderfile struct.  It also knows how to
tell if the Builderfile is valid (openable) or nat.
*/
type Parser struct {
	filename string
	log.Log
	dclient dclient.Dclient
}

/*
NewParser returns an initialized Parser.  Not currently necessary, as no
default values are assigned to a new Parser, but useful to have in case we need
to change this.
*/
func NewParser(filename string, logger log.Log) (*Parser, error) {
	fmt.Printf("logger here: %+v\n", logger)
	client, err := dclient.NewDclient(logger)
	if err != nil {
		return nil, err
	}

	return &Parser{
		Log:      logger,
		filename: filename,
		dclient:  *client,
	}, nil
}
