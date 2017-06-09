package media

import (
	"github.com/bogem/id3v2"
	"github.com/tcolgate/mp3"
	"net/http"
	"os"
	"path/filepath"
	"html/template"
	"strings"
	"time"
)

const TAG_TRACK = "TRCK"

type MediaContent struct {
	Name string
	Path string
	Duration string
}

type PageContext struct {
	MediaList map[string]MediaContent
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	pctx := &PageContext{
		MediaList: make(map[string]MediaContent),
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
			duration := ""
			if strings.HasSuffix(fname, ".mp3") {
				// parse the header
				func (){
					tag, err := id3v2.Open(fpath, id3v2.Options{Parse: true})
					if err != nil {
						os.Stderr.WriteString("ERROR: couldn't open or parse " + fpath + ": " + err.Error() + "\n")
					} else {
						defer tag.Close()

						if tag.Count() > 0 {
							track := strings.TrimSpace(tag.GetTextFrame(TAG_TRACK).Text)
							if track != "" {
								track += ": "
							}
							artist := strings.TrimSpace(tag.Artist())
							if artist != "" {
								artist += " - "
							}
							title := strings.TrimSpace(tag.Title())
							if title == "" && artist == "" {
								title = fname
							} else {
								title = track + artist + title
							}
							fname = title
						}
					}
				}()

				// parse the file content
				mp3f, err := os.Open(fpath)
				if err != nil {
					os.Stderr.WriteString("ERROR:  Couldn't open file " + fname + " for reading: %v" + err.Error() + "\n")
				}
				func () {
					skipped := 0
					defer func() {
						mp3f.Close()
					}()
					decoder := mp3.NewDecoder(mp3f)
					var trackDuration time.Duration
					var frame mp3.Frame
					for {
						if err := decoder.Decode(&frame, &skipped); err != nil {
							os.Stdout.WriteString("INFO: decode error: " + err.Error() + "\n")
							break
						}
						trackDuration += frame.Duration()
					}
					trackDuration = time.Duration(trackDuration.Seconds()) * time.Second
					duration = trackDuration.String()
				}()

			}
			content := MediaContent{
				Name: fname,
				Path: fpath,
				Duration: duration,
			}
			pctx.MediaList[fname] = content
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
