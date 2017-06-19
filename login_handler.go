package media

import (
	"net/http"
	"github.com.old/goware/emailx"
	"os"
	"html/template"
)

const PATH_LOGIN = "login"

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	pctx := &PageContext{}
	r.ParseForm()

	emailAddress := r.Form.Get("email_address")
	if emailAddress == "" {
		pctx.Error = "No email address provided."
	} else if err := emailx.Validate(r.Form.Get("email")); err != nil {
		pctx.Error = err.Error()
	} else {
		// valid email address, generate and send token email
		err := sendTokenEmail(r.URL, emailAddress)
		if err != nil {
			pctx.Error = "Couldn't send login token email. " + err.Error()
		} else {
			pctx.Info = "Sent login token email.  Please log in using the link in the email message."
		}
	}

	// create the template
	loginTemplate, err := template.New("login").ParseFiles("templates/login.html")
	if err != nil {
		os.Stderr.WriteString("Failed to load login template: " + err.Error() + "\n")
	}

	// pass the page context as the parameters to be processed by the template
	err = loginTemplate.ExecuteTemplate(w, "login", pctx)
	if err != nil {
		os.Stderr.WriteString("Failed to execute login template: " + err.Error() + "\n")
	}

}
