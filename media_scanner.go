package media

import (
	"os"
	"path/filepath"
	"strings"
)

func scanMedia(root string, media map[string]*MediaContent) {
	// scan the media
	mediaPath, err := os.Open(root)
	if err != nil {
		os.Stderr.WriteString("ERROR: failed to open media path " + root + ": " + err.Error() + "\n")
	}
	files, err := mediaPath.Readdir(0)
	if err != nil {
		os.Stderr.WriteString("ERROR: listing media path " + root + ": " + err.Error() + "\n")
	}
	for _, file := range files {
		if !file.IsDir() {
			fname := file.Name()
			fpath := filepath.Join(root, fname)
			if strings.HasSuffix(fname, ".mp3") {
				content := parseMp3File(fname, fpath)
				media[fname] = content
			}
		} else {
			scanMedia(filepath.Join(root, file.Name()), media)
		}
	}


}
