package conf

// Config is the global config
var Config Conf

// Conf is used for storing data retrieved from environmental variables.
type Conf struct {
	LogLevel  string `envconfig:"LOG_LEVEL"`
	LogFormat string `envconfig:"LOG_FORMAT"`
	APIToken  string `envconfig:"API_TOKEN"`
	SkipPush  bool   `envconfig:"SKIP_PUSH"`

	// docker registry credentials
	CfgUn    string `envconfig:"CFG_UN"`
	CfgPass  string `envconfig:"CFG_PASS"`
	CfgEmail string `envconfig:"CFG_EMAIL"`
}
