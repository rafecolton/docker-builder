package builderfile

/*
Builderfile is a struct representation of what is expected to be inside a
Builderfile.
*/
type Builderfile struct {
	Docker     `toml:"docker"`
	Containers map[string]ContainerSection
}

/*
Docker is a struct representation of the "docker" section of a Builderfile.
*/
type Docker struct {
	BuildOpts []string `toml:"build_opts"`
	TagOpts   []string `toml:"tag_opts"`
}

/*
ContainerSection is a struct representation of an individual member of the  "containers"
section of a Builderfile. Each of these sections defines a docker container to
be built and other related options.
*/
type ContainerSection struct {
	Dockerfile string
	Included   []string
	Excluded   []string
	Registry   string
	Project    string
	Tags       []string
}
