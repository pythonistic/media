package media

import (
	"net/http"
	"os"
	"html/template"
	"fmt"
	"time"
	"encoding/hex"
)

const COOKIE_MEDIA_USER = "media_user"
const COOKIE_KEY = "bb92c88dac984f677c1019d46c974fa8"
const TWO_YEARS = time.Duration(2 * 365 * 24) * time.Hour

var cookieKey = []byte(COOKIE_KEY)
var artists *Artists
var playlists *Playlists
var users *Users

func SetArtists(a *Artists) {
	artists = a
}

func SetPlaylists(p *Playlists) {
	playlists = p
}

func SetUsers(u *Users) {
	users = u
}

func GetUserEmail(r *http.Request) (email string){
	for _, cookie := range r.Cookies() {
		if cookie.Name == COOKIE_MEDIA_USER {
			nonce, err := hex.DecodeString(cookie.Value[0:12])
			if err != nil {
				os.Stderr.WriteString("Failed to decode cookie nonce: " + err.Error() + "\n")
			}
			cryptext, err := hex.DecodeString(cookie.Value[12:])
			if err != nil {
				os.Stderr.WriteString("Failed to decode cookie cryptext: " + err.Error() + "\n")
			}
			email = string(Decrypt(nonce, cryptext))
			break
		}
	}

	return
}

func SetUserEmail(w http.ResponseWriter, r *http.Request, email string) {
	nonce, cryptext := Encrypt([]byte(email))
	cookie := &http.Cookie{
		Value: fmt.Sprintf("%x%x", nonce, cryptext),
		Name: COOKIE_MEDIA_USER,
		Expires: time.Now().Add(TWO_YEARS),
		Secure: r.URL.Scheme == "https",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	pctx := &PageContext{
		Artists: artists,
	}

	// get the user email from the cookie
	email := GetUserEmail(r)

	if email != "" {
		user, ok := users.GetByEmail(email)
		if !ok {
			os.Stderr.WriteString("ERROR: couldn't find user by email " + email + "\n")
		}
		pctx.Playlists = playlists.GetForUser(user)
		pctx.User = user
	}

	// create the template
	indexTemplate, err := template.New("index").ParseFiles("templates/index.html")
	if err != nil {
		os.Stderr.WriteString("Failed to load index template: " + err.Error() + "\n")
	}

	// add functions to the templates
	//funcMap := template.FuncMap{
	//	"SomeFunction": func() string {
	//		var htmlFragment bytes.Buffer
	//
	//		return htmlFragment.String()
	//	},
	//}

	// pass the page context as the parameters to be processed by the template
	err = indexTemplate.ExecuteTemplate(w, "index", pctx)
	if err != nil {
		os.Stderr.WriteString("Failed to execute index template: " + err.Error() + "\n")
	}
}
