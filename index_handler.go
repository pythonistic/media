package media

import (
	"net/http"
	"os"
	"html/template"
)


func IndexHandler(w http.ResponseWriter, r *http.Request) {
	pctx := &PageContext{
		MediaList: make(map[string]*MediaContent),
	}

	scanMedia(FSPATH_MEDIA, pctx.MediaList)

	// render the template
	indexTemplate, err := template.New("index").ParseFiles("templates/index.html")
	if err != nil {
		os.Stderr.WriteString("Failed to load index template: " + err.Error() + "\n")
	}

	// pass the page context as the parameters to be processed by the template
	err = indexTemplate.ExecuteTemplate(w, "index", pctx)
	if err != nil {
		os.Stderr.WriteString("Failed to execute index template: " + err.Error() + "\n")
	}
}
