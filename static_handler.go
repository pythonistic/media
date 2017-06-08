package media

import (
	"net/http"
	"strings"
	"os"
)

const FSPATH_STATIC = "static"
const PATH_STATIC = "static"
const PATH_IMAGES = "images"
const PATH_FAVICON = "/favicon.ico"
const PATH_APPLE_TOUCH = "/apple-touch-icon.png"
const PATH_APPLE_PRECOMP = "/apple-touch-icon-precomposed.png"
const PATH_ANDROID_HIRES = "/icon-hires.png"
const PATH_ANDROID_NORMAL = "/icon-normal.png"

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	// handle icons - change the request path to static/images directory
	reqPath := r.URL.Path
	if reqPath == PATH_APPLE_PRECOMP || reqPath == PATH_APPLE_TOUCH ||
		reqPath == PATH_ANDROID_HIRES || reqPath == PATH_ANDROID_NORMAL ||
		reqPath == PATH_FAVICON {
		reqPath = PATH_STATIC + "/images" + reqPath
	}

	// strip out the leading static path and slash from the URL
	urlPath := reqPath[len(PATH_STATIC) + 1:]

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
	filePath := FSPATH_STATIC + "/" + urlPath

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
