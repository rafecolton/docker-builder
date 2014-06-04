package builderfile

/*
Builderfile is a struct representation of what is expected to be inside a
Builderfile.
*/
type Builderfile struct {
	Docker     `toml:"docker"`
	Containers map[string]ContainerSection `toml:"containers"`
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
	Name       string   `toml:"-"`
	Dockerfile string   `toml:"Dockerfile"`
	Included   []string `toml:"included"`
	Excluded   []string `toml:"excluded"`
	Registry   string   `toml:"registry"`
	Project    string   `toml:"project"`
	Tags       []string `toml:"tags"`
	SkipPush   bool     `toml:"skip_push"`
}

/*
Clean tidies up the structure of the Builderfile struct slightly by replacing
some occurrences of nil arrays with empty arrays []string{}.
*/
func (file *Builderfile) Clean() {
	if file.Docker.BuildOpts == nil {
		file.Docker.BuildOpts = []string{}
	}

	if file.Docker.TagOpts == nil {
		file.Docker.TagOpts = []string{}
	}
}
