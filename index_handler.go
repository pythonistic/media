package media

import (
	"net/http"
	"os"
	"html/template"
)

var artists *Artists

func SetArtists(a *Artists) {
	artists = a
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	pctx := &PageContext{
		Artists: artists,
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
