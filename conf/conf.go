package conf

// Config is the global config for docker-builder
var Config Conf

// Conf is used for storing data retrieved from environmental variables.
type Conf struct {
	Port      int
	LogLevel  string
	LogFormat string
	APIToken  string
	SkipPush  bool
	// for basic auth
	Username string
	Password string
	// for travis auth
	TravisToken  string
	GitHubSecret string
	NoTravis     bool
	NoGitHub     bool
	SleepTime    int
}
