package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

type SecretsVerifierMiddleware struct {
	handler http.Handler
}

func NewSecretsVerifierMiddleware(
	h http.Handler) *SecretsVerifierMiddleware {
	return &SecretsVerifierMiddleware{h}
}

func (v *SecretsVerifierMiddleware) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request) {

	fmt.Println("[Start]ServeHTTP")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	sv, err := slack.NewSecretsVerifier(r.Header,
		os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("before hander.ServeHTTP")
	v.handler.ServeHTTP(w, r)

	fmt.Println("[END]ServeHTTP")
}
