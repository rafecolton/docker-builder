package kamino

/*
A CacheOption is one of the valid cache option constants.
*/
type CacheOption string

const (
	// Create -  use cache if already created, create cache if not present.
	Create CacheOption = "create"

	// Force - use cache if already created, fail if cache not present.
	Force CacheOption = "force"

	// IfAvailable - use cache if already created, otherwise create a uniquely named directory.
	IfAvailable CacheOption = "if_available"

	// No - do not use cache, create a uniquely named directory.
	No CacheOption = "no"
)

/*
A Genome is the genetic options for a clone.  Here are the requirements for a valid genome:

	* Account
		- purpose: GitHub account
		- required: true
	* Cache
		- purpose: whether or not to use the cached / previously cloned version of the repo
		- required: false
		- default: No
	* Depth
		- purpose: git clone `--depth` option
		- required: false
		- default: "" (the whole repo will be cloned)
		- validation: must be empty string or parsable as a base 10 integer
	* Ref
		- purpose: the git SHA to check out in the cloned repo
		- required: true
	* Repo
		- purpose: GitHub repo
		- required: true
	* APIToken
		- purpose: GitHub API token for private repos
		- required: false (functionally required if your repo is private)
		- default: (not sent with request if empty)

*/
type Genome struct {
	APIToken  string
	Account   string
	Depth     string
	Recursive bool
	Ref       string
	Repo      string
	UseCache  CacheOption
}

/*
IsValid indicates whether or not the cache option is valid.
*/
func (opt CacheOption) IsValid() bool {
	return opt == Create || opt == Force || opt == IfAvailable || opt == No
}
