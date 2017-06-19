package media

import (
	"net/http"
	"strings"
)

const PATH_TOKEN = "token"

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	// read the token from the request path
	baseUrl := r.URL.String()[0:strings.LastIndex(r.URL.String(), "/token")]
	code := r.URL.Path[strings.LastIndex(r.URL.Path, "/"):]

	token, err := processEmailToken(code)
	if token == nil || err == nil {
		// show the login page because the token wasn't found
		// TODO consider tossing an error message in here someplace to tell the user about the failed login
		http.Redirect(w, r, baseUrl, http.StatusTemporaryRedirect)
		return
	}

	// write the cookie
	SetUserEmail(w, r, token.Email)

	// redirect the user to the main page
	http.Redirect(w, r, baseUrl, http.StatusTemporaryRedirect)
}