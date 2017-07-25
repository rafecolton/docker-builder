package unitconfig

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/winchman/builder-core/filecheck"
	"gopkg.in/yaml.v2"
)

// Encoding is a constant that represents the config format of a file
// that is expected to contain an encoded unit-config
type Encoding int

const (
	// TOML is for files encoded with TOML (or ending in .toml)
	TOML Encoding = iota

	// JSON is for files encoded in JSON (or ending in .json)
	JSON

	// YAML is for files encoded in YAML (or ending in .yml or .yaml)
	YAML

	unknown
)

// ReadFromFile accepts an unsanitized path to a Bobfile and returns the
// decoded UnitConfig.  The value of encodings is a list of possible file
// encodings.  If none is provided, the encoding will be inferred from the file
// extension as described in the documentation for the UnitConfigEncoding
// constants.  If no encoding can be inferred, TOML will be used as the default.
func ReadFromFile(path string, encodings ...Encoding) (*UnitConfig, error) {
	// check file
	if path == "" {
		return nil, errors.New("no file path provided")
	}

	top := filepath.Dir(path)
	filename := filepath.Base(path)
	opts := filecheck.NewTrustedFilePathOptions{
		File: filename,
		Top:  top,
	}

	var file *filecheck.TrustedFilePath
	var err error
	if file, err = filecheck.NewTrustedFilePath(opts); err != nil {
		return nil, err
	}
	if file.Sanitize(); file.State != filecheck.OK {
		return nil, errors.New("provided file is unchecked or unsanitary")
	}

	// file is sanitary

	contents, err := ioutil.ReadFile(file.FullPath())
	if err != nil {
		return nil, err
	}

	// attempt to infer encoding by file extension
	if encodings == nil || len(encodings) == 0 {
		inferred := inferredEncoding(file)
		if inferred == unknown {
			inferred = TOML
		}
		encodings = []Encoding{inferred}
	}

	// return first validly-decoded UnitConfig
	for _, encoding := range encodings {
		switch encoding {
		case TOML:
			ret := &UnitConfig{}
			if _, err = toml.Decode(string(contents), &ret); err == nil {
				return ret, nil
			}
		case JSON:
			decoder := json.NewDecoder(bytes.NewReader(contents))
			ret := &UnitConfig{}
			if err := decoder.Decode(ret); err == nil {
				return ret, nil
			}
		case YAML:
			ret := &UnitConfig{}
			if err = yaml.Unmarshal(contents, &ret); err == nil {
				return ret, nil
			}
		}
	}

	if err == nil {
		err = errors.New("unable to decode file contents")
	}

	return nil, err
}

func inferredEncoding(file *filecheck.TrustedFilePath) Encoding {
	ext := strings.ToLower(filepath.Ext(file.File()))
	switch ext {
	case "json":
		return JSON
	case "toml":
		return TOML
	case "yaml", "yml":
		return YAML
	default:
		return unknown
	}
}
