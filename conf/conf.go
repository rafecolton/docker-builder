package conf

import (
	"math/rand"
	"strings"
	"time"
)

// Config is the global config for docker-builder
var Config Conf

// Secret is a type that wraps string, causing cli to print out "*" instead of
// the actual value
type Secret string

// Conf is used for storing data retrieved from environmental variables.
type Conf struct {
	Port      int
	LogLevel  string
	LogFormat string
	APIToken  string
	SkipPush  bool
	Squash    bool

	// for basic auth
	Username string
	Password string

	// for travis auth
	TravisToken string
	NoTravis    bool

	// for github auth
	GitHubSecret string
	NoGitHub     bool

	// docker registry credentials
	CfgUn    Secret
	CfgPass  Secret
	CfgEmail Secret
}

func (s Secret) String() string {
	var length int
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	positive := random.Intn(2)
	offset := random.Intn(10)

	if positive == 0 {
		length = len(s) + offset
	} else {
		length = len(s) - offset
		if length < 1 {
			length = 1
		}
	}
	return strings.Repeat("*", length)
}
