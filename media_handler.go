package media

import (
	"net/http"
	"strings"
	"os"
)

const PATH_MEDIA = "media"
const FSPATH_MEDIA = "media"

func MediaHandler(w http.ResponseWriter, r *http.Request) {
	// strip out the leading media path and slash from the URL
	urlPath := r.URL.Path[len(PATH_MEDIA) + 1:]

	// strip directory traversals
	for strings.Contains(urlPath, "..") {
		urlPath = strings.Replace(urlPath, "..", "", -1)
	}

	// strip out any leading absolute path
	for strings.HasPrefix(urlPath, "/") {
		urlPath = urlPath[1:]
	}

	// prepend the urlPath with the real path
	// NOTE: Not Windows NTFS safe; sorry Windows users!
	filePath := FSPATH_MEDIA + "/" + urlPath

	// verify the path exists before we try to open it
	fileinfo, err := os.Stat(filePath)
	if err != nil || fileinfo.IsDir() || !fileinfo.Mode().IsRegular() {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// else, use the stdlib content handler (for MIME type support, range calls, etc.)
	parts := strings.Split(filePath, "/")
	fs := MediaFileSystem{}
	seeker, err := fs.Open(filePath)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		os.Stderr.WriteString("ERROR: Forbidden open " + filePath + ": " + err.Error() + "\n")
	}
	http.ServeContent(w, r, parts[len(parts) - 1], fileinfo.ModTime(), seeker)
}
