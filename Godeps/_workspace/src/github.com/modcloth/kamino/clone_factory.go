// Package kamino makes it super easy to clone things from GitHub.
package kamino

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
)

/*
A CloneFactory generates your clones for you.  Create a clone factory with NewCloneFactory().
*/
type CloneFactory struct {
	workdir string
}

/*
NewCloneFactory creates a new CloneFactory, ready to do some cloning for you
(into its specified workdir).
*/
func NewCloneFactory(workdir string) (*CloneFactory, error) {

	if err := os.MkdirAll(workdir, 0755); err != nil {
		return nil, err
	}

	return &CloneFactory{
		workdir: workdir,
	}, nil
}

/*
Clone clones the repo as specified by the genome parameters.
*/
func (factory *CloneFactory) Clone(g *Genome) (path string, err error) {
	if err = ValidateGenome(g); err != nil {
		return "", err
	}

	creator := &clone{
		Genome:  g,
		workdir: factory.workdir,
	}

	Logger.WithFields(logrus.Fields{
		"account":            g.Account,
		"cache_method":       g.UseCache,
		"depth":              g.Depth,
		"ref":                g.Ref,
		"repo":               g.Repo,
		"api_token_provided": g.APIToken != "",
	}).Debug("requesting clone")

	switch g.UseCache {
	case No:
		return creator.cloneNoCache()
	case Create:
		return creator.cloneCreateCache()
	case Force:
		return creator.cloneForceCache()
	case IfAvailable:
		return creator.cloneCacheIfAvailable()
	default:
		return creator.cloneNoCache()
	}
}

/*
ValidateGenome validates a genome and returns a non-nil error if the genome is invalid.
*/
func ValidateGenome(g *Genome) error {
	if g.Depth != "" {
		if _, err := strconv.Atoi(g.Depth); err != nil {
			return fmt.Errorf("%q is not a valid clone depth", g.Depth)
		}
	}

	if !g.UseCache.IsValid() {
		return fmt.Errorf("%q is not a valid cache option", g.UseCache)
	}

	if g.Account == "" {
		return errors.New("account must be provided")
	}

	if g.Ref == "" {
		return errors.New("ref must be provided")
	}

	if g.Repo == "" {
		return errors.New("repo must be provided")
	}

	return nil
}
