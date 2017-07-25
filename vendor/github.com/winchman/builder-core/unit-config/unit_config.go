package unitconfig

/*
UnitConfig is a struct representation of what is expected to be inside a
Builderfile for a single build/tag/push sequence.
*/
type UnitConfig struct {
	Version          int                 `toml:"version" json:"version" yaml:"version"`
	Docker           Docker              `toml:"docker" json:"docker" yaml:"docker"`
	ContainerArr     []*ContainerSection `toml:"container" json:"container" yaml:"container"`
	ContainerGlobals *ContainerSection   `toml:"container_globals" json:"container_globals" yaml:"container_globals"`
}

/*
Docker is a struct representation of the "docker" section of a Builderfile.
*/
type Docker struct {
	BuildOpts []string `toml:"build_opts" json:"build_opts" yaml:"build_opts"`
	TagOpts   []string `toml:"tag_opts" json:"tag_opts" yaml:"tag_opts"`
}

/*
ContainerSection is a struct representation of an individual member of the  "containers"
section of a Builderfile. Each of these sections defines a docker container to
be built and other related options.
*/
type ContainerSection struct {
	Name       string   `toml:"name" json:"name" yaml:"name"`
	Dockerfile string   `toml:"Dockerfile" json:"Dockerfile" yaml:"Dockerfile"`
	Registry   string   `toml:"registry" json:"registry" yaml:"registry"`
	Project    string   `toml:"project" json:"project" yaml:"project"`
	Tags       []string `toml:"tags" json:"tags" yaml:"tags"`
	SkipPush   bool     `toml:"skip_push" json:"skip_push" yaml:"skip_push"`
	CfgUn      string   `toml:"dockercfg_un" json:"dockercfg_un" yaml:"dockercfg_un"`
	CfgPass    string   `toml:"dockercfg_pass" json:"dockercfg_pass" yaml:"dockercfg_pass"`
	CfgEmail   string   `toml:"dockercfg_email" json:"dockercfg_email" yaml:"dockercfg_email"`
}

// ConfigGlobals are the global values that may be set with SetGlobals
type ConfigGlobals struct {
	SkipPush bool
	CfgUn    string
	CfgPass  string
	CfgEmail string
}

// SetGlobals allows some global values to be set (i.e. if pulled from the env in the calling process)
func (config *UnitConfig) SetGlobals(globals ConfigGlobals) {
	if config.ContainerGlobals == nil {
		config.ContainerGlobals = &ContainerSection{}
	}

	if !config.ContainerGlobals.SkipPush {
		config.ContainerGlobals.SkipPush = globals.SkipPush
	}
	if config.ContainerGlobals.CfgUn == "" {
		config.ContainerGlobals.CfgUn = globals.CfgUn
	}
	if config.ContainerGlobals.CfgPass == "" {
		config.ContainerGlobals.CfgPass = globals.CfgPass
	}
	if config.ContainerGlobals.CfgEmail == "" {
		config.ContainerGlobals.CfgEmail = globals.CfgEmail
	}

}
