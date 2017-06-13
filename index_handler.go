package media

import (
	"net/http"
	"os"
	"html/template"
)


func IndexHandler(w http.ResponseWriter, r *http.Request) {
	pctx := &PageContext{
		Artists: make(Artists, 0),
	}

	scanMedia(FSPATH_MEDIA, pctx.Artists)

	// create the template
	indexTemplate, err := template.New("index").ParseFiles("templates/index.html")
	if err != nil {
		os.Stderr.WriteString("Failed to load index template: " + err.Error() + "\n")
	}

	// add functions to the templates
	//funcMap := template.FuncMap{
	//	"MatchedDomain": func() string {
	//		var htmlFragment bytes.Buffer
	//		matchedDomain := ""
	//		zuulHosts := strings.Split(_auth_config.Config("handler.hosts"), ",")
	//
	//		for idx, domain := range zuulHosts {
	//			if strings.Contains(pctx.Request.Host, domain) {
	//				// copy the request URL so we can safely modify it
	//				var ssoUrl *url.URL = pctx.Request.URL
	//				ssoUrl.Path = "signoutSSO"
	//				ssoUrl.Host = domain
	//
	//				logoutDomainTemplate.ExecuteTemplate(&htmlFragment, "logout_domain", struct {
	//					Index         int
	//					MatchedDomain string
	//					MatchedSsoUrl string
	//				}{
	//					idx, matchedDomain, ssoUrl.String(),
	//				})
	//			}
	//		}
	//		return htmlFragment.String()
	//	},
	//}

	// pass the page context as the parameters to be processed by the template
	err = indexTemplate.ExecuteTemplate(w, "index", pctx)
	if err != nil {
		os.Stderr.WriteString("Failed to execute index template: " + err.Error() + "\n")
	}
}
