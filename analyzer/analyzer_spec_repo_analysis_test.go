package analyzer_test

type SpecRepoAnalysis struct {
	remotes           string
	dockerfilePresent bool
	isGitRepo         bool
	repoBasename      string
}

func (sra *SpecRepoAnalysis) GitRemotes() string {
	return sra.remotes
}

func (sra *SpecRepoAnalysis) DockerfilePresent() bool {
	return sra.dockerfilePresent
}

func (sra *SpecRepoAnalysis) IsGitRepo() bool {
	return sra.isGitRepo
}

func (sra *SpecRepoAnalysis) RepoBasename() string {
	return sra.repoBasename
}
