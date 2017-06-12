package media

import (
	"github.com/bogem/id3v2"
	"os"
	"strings"
	"github.com/tcolgate/mp3"
	"time"
)

func parseMp3File(fname string, fpath string) (media *MediaContent) {
	media = &MediaContent{
		Filename: fname,
		Path: fpath,
	}

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

				media.Title = title
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
		media.Duration = trackDuration.String()
	}()

	return
}