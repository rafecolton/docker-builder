package vauth

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

/*
TravisCI returns a Handler that authenticates via Travis's Authorization for
Webhooks scheme (http://docs.travis-ci.com/user/notifications/#Authorization-for-Webhooks)

Writes a http.StatusUnauthorized if authentication fails
*/
func TravisCI(token string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		providedAuth := req.Header.Get("Authorization")

		travisRepoSlug := req.Header.Get("Travis-Repo-Slug")
		calculatedAuth := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s%s", travisRepoSlug, token))))

		if !SecureCompare(providedAuth, calculatedAuth) {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		}
	}
}
