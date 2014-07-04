package conf

var Config Conf

/*
Config is used for storing data retrieved from environmental variables.
*/
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
}
