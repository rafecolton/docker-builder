package vauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
)

/*
GitHub returns a Handler that authenticates via GitHub's Authorization for
Webhooks scheme (https://developer.github.com/webhooks/securing/#validating-payloads-from-github)

Writes a http.StatusUnauthorized if authentication fails
*/
func GitHub(secret string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		requestSignature := req.Header.Get("X-Hub-Signature")

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		}

		req.Body = ioutil.NopCloser(bytes.NewReader(body))

		mac := hmac.New(sha1.New, []byte(secret))
		mac.Reset()
		mac.Write(body)
		calculatedSignature := fmt.Sprintf("sha1=%x", mac.Sum(nil))

		if !SecureCompare(requestSignature, calculatedSignature) {
			http.Error(res, "Not Authorized", http.StatusUnauthorized)
		}
	}
}
