package media

import (
	"net/http"
	"os"
	"path/filepath"
	"html/template"
)

type PageContext struct {
	MediaList map[string]string
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	pctx := &PageContext{
		MediaList: make(map[string]string),
	}

	// scan the media
	media, err := os.Open(FSPATH_MEDIA)
	if err != nil {
		os.Stderr.WriteString("ERROR: failed to open media path " + FSPATH_MEDIA + ": " + err.Error() + "\n")
	}
	files, err := media.Readdir(0)
	if err != nil {
		os.Stderr.WriteString("ERROR: listing media path " + FSPATH_MEDIA + ": " + err.Error() + "\n")
	}
	for _, file := range files {
		if !file.IsDir() {
			fname := file.Name()
			fpath := filepath.Join(FSPATH_MEDIA, fname)
			pctx.MediaList[fname] = fpath
		}
	}

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
