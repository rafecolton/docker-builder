package builderfile

type Builderfile struct {
	Docker     `toml:"docker"`
	Containers map[string]ContainerSection
}

type Docker struct {
	BuildOpts string `toml:"build_opts"`
}

type ContainerSection struct {
	Dockerfile string
	Included   []string
	Excluded   []string
	Registry   string
	Project    string
	Tags       []string
}
